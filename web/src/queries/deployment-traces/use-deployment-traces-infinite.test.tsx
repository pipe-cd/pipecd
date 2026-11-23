// Copyright 2026 The PipeCD Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

import { act, renderHook } from "@testing-library/react-hooks";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { FC } from "react";
import { dummyDeploymentTrace } from "~/__fixtures__/dummy-deployment-trace";
import { getDeploymentTraces } from "~/api/deploymentTraces";
import { useGetDeploymentTracesInfinite } from "./use-deployment-traces-infinite";

jest.mock("~/api/deploymentTraces");

const mockGetDeploymentTraces = jest.mocked(getDeploymentTraces);

it.each([50, 10, 0])(
  "loads the next page after an initial page of %i traces",
  async (traceCount) => {
    const queryClient = new QueryClient({
      defaultOptions: { queries: { retry: false, cacheTime: 0 } },
    });
    const wrapper: FC = ({ children }) => (
      <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
    );
    mockGetDeploymentTraces
      .mockResolvedValueOnce({
        tracesList: Array(traceCount).fill(dummyDeploymentTrace),
        cursor: "next-cursor",
      })
      .mockResolvedValueOnce({
        tracesList: [dummyDeploymentTrace],
        cursor: "older-cursor",
      });

    const { result, waitFor } = renderHook(
      () => useGetDeploymentTracesInfinite({ commitHash: "commit-hash" }),
      { wrapper }
    );
    await waitFor(() => expect(result.current.isSuccess).toBe(true));
    expect(result.current.data?.pages[0].hasMore).toBe(traceCount === 50);

    const firstRequest = mockGetDeploymentTraces.mock.calls[0][0];
    act(() => result.current.fetchNextPage());
    await waitFor(() => expect(result.current.data?.pages).toHaveLength(2));

    expect(mockGetDeploymentTraces).toHaveBeenLastCalledWith({
      options: { commitHash: "commit-hash" },
      pageSize: 30,
      cursor: "next-cursor",
      pageMinUpdatedAt:
        firstRequest.pageMinUpdatedAt - (traceCount === 50 ? 0 : 2592000),
    });
    expect(result.current.data?.pages[1].tracesList).toEqual([
      dummyDeploymentTrace,
    ]);
    queryClient.clear();
  }
);
