fetch("sentry-cfg.json").then(function (response) {
  response.json().then(function (result) {
    sentry_default = {
      integrations(integrations) {
        return integrations.filter(integration => integration.name !== 'Breadcrumbs');
      }
    }
    if (result.dsn) {
      Sentry.init(Object.assign(result, sentry_default));
    }
  });
});

function extractKey(k) {
  const hashIdx = "_**[]"
  if (k.indexOf(hashIdx) < 0) {
    return k
  }
  return k.split(hashIdx)[0].trim()
}

function getRequestBody(str) {
  try {
    return JSON.stringify(JSON.parse(str), null, 2);
  } catch (e) {
    return str;
  }
}

function curl(request, data) {
  let curlified = []
  let type = ""
  let headers = request.headers
  curlified.push(`curl -X ${data.method}`)
  curlified.push(`"${data.url}"`)
  if (headers && Object.entries(data.headers).length) {
    // #Object.entries(data.headers)
    for (let p of request.headers.entries()) {
      let [h, v] = p
      type = v
      curlified.push(`-H "${h}: ${v}"`)
    }
  }

  if (data.body) {
    if (type.includes("multipart/form-data") && request.method === "POST") {
      for (let [k, v] of data.body.entries()) {
        let extractedKey = extractKey(k)
        if (v instanceof File) {
          curlified.push(`-F "${extractedKey}=@${v.name}${v.type ? `;type=${v.type}` : ""}"`)
        } else {
          curlified.push(`-F "${extractedKey}=${v}"`)
        }
      }
    } else {
      curlified.push(`-d ${JSON.stringify(data.body).replace(/\\n/g, "")}`)
    }
  }
  return curlified.join(" \\\n")
}

var origFetch = fetch;
fetch = function (url, data = {}) {
  req = new Request(url, data);
  return origFetch(req).then(function (response) {
    if (!url.includes("sentry")) {
      res = response.clone();
      if (!res.ok) {
        res.text().then(function (result) {
          request_curl = curl(req, data);
          response_body = result;
          var contentType = res.headers.get("content-type");
          if (contentType && contentType.includes("application/json")) {
            response_body = JSON.stringify(JSON.parse(result), null, 2);
          }
          response_url = new URL(res.url);
          response_headers = '';
          for (let [key, value] of res.headers) {
            response_headers += `${key}: ${value}\n`;
          }
          Sentry.withScope(function (scope) {
            scope.setTag("response.code", res.status);
            if (req.headers.get("accept") == "application/json" && data.body) {
              scope.setExtra("Request Body", getRequestBody(data.body));
            }
            scope.setExtra("Request CURL", curl(req, data));
            scope.setExtra("Response Headers", response_headers);
            scope.setExtra("Response Body", response_body);
            scope.addEventProcessor(sentry_event => {
              sentry_event.request.url = res.url;
              sentry_event.request.headers = Object.assign(sentry_event.request.headers, data.headers);
              return sentry_event;
            });
            scope.setFingerprint([response_url]);
            Sentry.captureMessage(`${response_url.pathname} ${res.status}\n\n${response_body}`, 'error');
          });
        })
      }
    }
    return response;
  });
};
