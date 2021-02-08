export let apiEndpoint = `${location.protocol}//${location.host}`;

if (process.env.NODE_ENV === "development") {
  apiEndpoint = `${apiEndpoint}/api`;
}

if (process.env.NODE_ENV === "test") {
  apiEndpoint = "https://test.pipecd.dev";
}
