type Query = Record<string, any>;

// TODO: proper types for object arguments

var xfetch = function (resource: string, init: RequestInit = {}) {
  if (["post", "put", "delete"].indexOf(init.method || "") !== -1) {
    init["headers"] = new Headers(init["headers"]);
    init["headers"].set("x-requested-by", "yarr");
  }
  return fetch(resource, init);
};
var api = function (method: string, endpoint: string, data?: object) {
  var headers = { "Content-Type": "application/json" };
  return xfetch(endpoint, {
    method: method,
    headers: headers,
    body: JSON.stringify(data),
  });
};

var json = function (res: Response) {
  return res.json();
};

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
      return api("post", "./api/feeds", data).then(json);
    },
    update(id: number, data: object) {
      return api("put", "./api/feeds/" + id, data);
    },
    delete(id: number) {
      return api("delete", "./api/feeds/" + id);
    },
    list_items(id: number) {
      return api("get", "./api/feeds/" + id + "/items").then(json);
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
      return api("post", "./api/folders", data).then(json);
    },
    update(id: number, data: object) {
      return api("put", "./api/folders/" + id, data);
    },
    delete(id: number) {
      return api("delete", "./api/folders/" + id);
    },
    list_items(id: number) {
      return api("get", "./api/folders/" + id + "/items").then(json);
    },
  },
  items: {
    get(id: number) {
      return api("get", "./api/items/" + id).then(json);
    },
    list(query: Query) {
      return api("get", "./api/items" + param(query)).then(json);
    },
    update(id: number, data: object) {
      return api("put", "./api/items/" + id, data);
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
      return api("put", "./api/settings", data);
    },
  },
  status() {
    return api("get", "./api/status").then(json);
  },
  upload_opml(form: HTMLFormElement) {
    return xfetch("./opml/import", {
      method: "post",
      body: new FormData(form),
    });
  },
  logout() {
    return api("post", "./logout");
  },
  crawl(url: string) {
    return api("get", "./page?url=" + encodeURIComponent(url)).then(json);
  },
};
