interface Window {
  app: any;
}

declare const process: {
  env: {
    NODE_ENV: string;
  };
};

declare function require(module: string): any;
