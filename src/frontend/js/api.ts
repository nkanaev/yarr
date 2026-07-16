type Query = Record<string, any>;

// TODO: proper types for object arguments

type ApiOptions = {
  json?: Record<string, any>;
  body?: BodyInit;
  query?: Record<string, string>;
};

function api(method: string, endpoint: string, opts: ApiOptions = {}) {
  const { json, body, query } = opts;

  let url = endpoint;

  const init: RequestInit = {};
  init.method = method;
  init.headers = new Headers();
  init.headers.set("x-requested-by", "yarr");

  if (query !== undefined) {
    url = url + "?" + new URLSearchParams(query).toString();
  }
  if (body !== undefined) {
    init.body = body;
  }
  if (json !== undefined) {
    init.headers.set("content-type", "application/json");
    init.body = JSON.stringify(json);
  }

  return fetch(url, init);
}

function json(res: Response) {
  return res.json();
}

function param(query: Query) {
  if (!query) return "";
  return (
    "?" +
    Object.keys(query)
      .map(function (key) {
        return encodeURIComponent(key) + "=" + encodeURIComponent(query[key]);
      })
      .join("&")
  );
}

export default {
  feeds: {
    list() {
      return api("get", "./api/feeds").then(json);
    },
    create(data: object) {
      return api("post", "./api/feeds", { json: data }).then(json);
    },
    update(id: number, data: object) {
      return api("put", `./api/feeds/${id}`, { json: data });
    },
    delete(id: number) {
      return api("delete", `./api/feeds/${id}`);
    },
    list_items(id: number) {
      return api("get", `./api/feeds/${id}/items`).then(json);
    },
    refresh() {
      return api("post", "./api/feeds/refresh");
    },
    list_errors() {
      return api("get", "./api/feeds/errors").then(json);
    },
  },
  folders: {
    list() {
      return api("get", "./api/folders").then(json);
    },
    create(data: object) {
      return api("post", "./api/folders", { json: data }).then(json);
    },
    update(id: number, data: object) {
      return api("put", `./api/folders/${id}`, { json: data });
    },
    delete(id: number) {
      return api("delete", `./api/folders/${id}`);
    },
    list_items(id: number) {
      return api("get", `./api/folders/${id}/items`).then(json);
    },
  },
  items: {
    get(id: number) {
      return api("get", `./api/items/${id}`).then(json);
    },
    list(query: Query) {
      return api("get", "./api/items", { query }).then(json);
    },
    update(id: number, data: object) {
      return api("put", `./api/items/${id}`, { json: data });
    },
    mark_read(query: Query) {
      return api("put", "./api/items" + param(query));
    },
  },
  settings: {
    get() {
      return api("get", "./api/settings").then(json);
    },
    update(data: object) {
      return api("put", "./api/settings", { json: data });
    },
  },
  status() {
    return api("get", "./api/status").then(json);
  },
  upload_opml(form: HTMLFormElement) {
    return api("post", "./opml/import", { body: new FormData(form) });
  },
  logout() {
    return api("post", "./logout");
  },
  crawl(url: string) {
    return api("get", "./page", { query: { url } }).then(json);
  },
};
