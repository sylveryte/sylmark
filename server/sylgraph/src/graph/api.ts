const baseUrl = "http://127.0.0.1:7462/v1/";

export async function httpGet<T>(url: RequestInfo | URL, init?: RequestInit) {
  const res = await fetch(baseUrl + url, {
    method: "GET",
    ...(init || {}),
  });
  return res.json() as T;
}

export async function httpPost<T>(
  url: RequestInfo | URL,
  body: any,
  init?: RequestInit,
) {
  const res = await fetch(baseUrl + url, {
    method: "POST",
    body: JSON.stringify(body),
    ...(init || {}),
  });
  return res.json() as T;
}
