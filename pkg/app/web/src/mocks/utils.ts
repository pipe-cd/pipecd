import { apiEndpoint } from "~/constants/api-endpoint";

export const createMask = (path: string): string =>
  `${apiEndpoint}/grpc.service.webservice.WebService${path}`;
