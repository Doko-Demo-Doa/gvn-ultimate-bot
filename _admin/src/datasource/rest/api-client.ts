// Same with hospital-api but rewritten with react-query
import qs from "qs";

/**
 * Important: Must be called first to initialize
 */
type InitParams = {
  baseUrl: string;
};

type RequestType<Req> = {
  method: string; //will define enum of method POST, GET etc
  url: string;
  headers?: Record<string, string>;
  body?: Req;
  multipart?: boolean;
};

const rawApi = async <Req>({
  body,
  headers,
  method,
  multipart,
  url,
}: RequestType<Req>) => {
  return await fetch(url, {
    body:
      body == null
        ? undefined
        : multipart
          ? (body as BodyInit)
          : JSON.stringify(body),
    headers,
    method,
  });
};

const _fetch = async <Req>(params: RequestType<Req>) => {
  const res = await rawApi(params);
  return res.json();
};

export class CustomApiClient {
  baseUrl = "";
  headers = {};

  init = (params: InitParams) => {
    this.baseUrl = params.baseUrl;
    this.headers = {
      "Content-Type": "application/json",
    };
  };

  setHeaders = (headers: Record<string, string>) => {
    this.headers = { ...this.headers, ...headers };
  };

  buildUrl = (url: string, useBaseUrl = true) => {
    if (!useBaseUrl || /^https?:\/\//.test(url)) {
      return url;
    }

    return `${this.baseUrl}${url}`;
  };

  get = async (url: string, query: Record<string, any>) => {
    return _fetch({
      headers: {
        ...this.headers,
        Accept: "application/ld+json",
        "Content-Type": "application/json",
      },
      method: "GET",
      url: `${this.buildUrl(url)}${qs.stringify(
        {
          ...query,
        },
        { addQueryPrefix: true, arrayFormat: "brackets", encode: false },
      )}`,
    });
  };

  // any is safe here
  post = async (url: string, body: any) => {
    return _fetch({
      body,
      headers: {
        ...this.headers,
        Accept: "application/ld+json",
        "Content-Type": "application/json",
      },
      method: "POST",
      url: this.buildUrl(url),
    });
  };

  postFormData = async (
    url: string,
    body: FormData,
    opts?: { useBaseUrl?: boolean },
  ) => {
    return rawApi({
      body,
      headers: {
        Accept: "*/*",
      },
      method: "POST",
      multipart: true,
      url: this.buildUrl(url, opts?.useBaseUrl),
    });
  };

  patch = async (url: string, body: any) => {
    return _fetch({
      body,
      headers: {
        ...this.headers,
        Accept: "application/ld+json",
        "Content-Type": "application/merge-patch+json",
      },
      method: "PATCH",
      url: this.buildUrl(url),
    });
  };

  delete = async (url: string) => {
    return _fetch({
      headers: {
        ...this.headers,
        Accept: "application/ld+json",
        "Content-Type": "application/json",
      },
      method: "DELETE",
      url: this.buildUrl(url),
    });
  };
}

export const customApiClient = new CustomApiClient();
