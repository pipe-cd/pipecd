import { apiEndpoint } from "../constants/api-endpoint";

export const createMask = (path: string): string =>
  `${apiEndpoint}/pipe.api.service.webservice.WebService${path}`;
