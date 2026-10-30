import { fireEvent, render, screen } from "~~/test-utils";
import { MemoryRouter } from "~~/test-utils";
import DeploymentTracePage from "./index";
import { dummyDeploymentTrace } from "~/__fixtures__/dummy-deployment-trace";
import { UI_TEXT_MORE, UI_TEXT_REFRESH } from "~/constants/ui-text";

let mockInView = false;
jest.mock("react-intersection-observer", () => ({
  useInView: () => [jest.fn(), mockInView],
}));

const mockFetchMoreDeploymentTraces = jest.fn();
const mockRefetchDeploymentTraces = jest.fn();

let mockQueryResult = {
  data: {
    pages: [
      {
        tracesList: [dummyDeploymentTrace],
        hasMore: true,
      },
    ],
  },
  isFetching: false,
  fetchNextPage: mockFetchMoreDeploymentTraces,
  refetch: mockRefetchDeploymentTraces,
  isSuccess: true,
};

jest.mock("~/queries/deployment-traces/use-deployment-traces-infinite", () => ({
  useGetDeploymentTracesInfinite: () => mockQueryResult,
}));

describe("DeploymentTracePage", () => {
  beforeEach(() => {
    jest.clearAllMocks();
    mockInView = false;
    mockQueryResult = {
      data: {
        pages: [
          {
            tracesList: [dummyDeploymentTrace],
            hasMore: true,
          },
        ],
      },
      isFetching: false,
      fetchNextPage: mockFetchMoreDeploymentTraces,
      refetch: mockRefetchDeploymentTraces,
      isSuccess: true,
    };
  });

  it("should render trace items when data is available", () => {
    render(
      <MemoryRouter>
        <DeploymentTracePage />
      </MemoryRouter>
    );

    expect(
      screen.getByText(dummyDeploymentTrace.trace?.title ?? "")
    ).toBeInTheDocument();
  });

  it("should render the More button when hasMore is true", () => {
    mockQueryResult.data.pages[0].hasMore = true;

    render(
      <MemoryRouter>
        <DeploymentTracePage />
      </MemoryRouter>
    );

    const moreButton = screen.getByRole("button", { name: UI_TEXT_MORE });
    expect(moreButton).toBeInTheDocument();
    expect(moreButton).toBeEnabled();

    fireEvent.click(moreButton);
    expect(mockFetchMoreDeploymentTraces).toHaveBeenCalledTimes(1);
  });

  it.each([10, 0])(
    "should allow loading older traces after a short page of %i traces",
    (traceCount) => {
      mockInView = true;
      mockQueryResult.data.pages[0] = {
        tracesList: Array.from({ length: traceCount }, (_, index) => ({
          ...dummyDeploymentTrace,
          trace: {
            ...dummyDeploymentTrace.trace!,
            id: `trace-${index}`,
          },
        })),
        hasMore: false,
      };

      const { rerender } = render(
        <MemoryRouter>
          <DeploymentTracePage />
        </MemoryRouter>
      );

      expect(mockFetchMoreDeploymentTraces).not.toHaveBeenCalled();
      const moreButton = screen.getByRole("button", { name: UI_TEXT_MORE });
      expect(moreButton).toBeEnabled();
      fireEvent.click(moreButton);
      expect(mockFetchMoreDeploymentTraces).toHaveBeenCalledTimes(1);

      mockQueryResult = { ...mockQueryResult, isFetching: true };
      rerender(
        <MemoryRouter>
          <DeploymentTracePage />
        </MemoryRouter>
      );
      expect(moreButton).toBeDisabled();
      fireEvent.click(moreButton);
      expect(mockFetchMoreDeploymentTraces).toHaveBeenCalledTimes(1);

      mockQueryResult = { ...mockQueryResult, isFetching: false };
      rerender(
        <MemoryRouter>
          <DeploymentTracePage />
        </MemoryRouter>
      );
      expect(moreButton).toBeEnabled();
      fireEvent.click(moreButton);
      expect(mockFetchMoreDeploymentTraces).toHaveBeenCalledTimes(2);
    }
  );

  it.each([true, false])(
    "should hide the More button before a successful query (isFetching: %s)",
    (isFetching) => {
      mockQueryResult.isSuccess = false;
      mockQueryResult.isFetching = isFetching;

      render(
        <MemoryRouter>
          <DeploymentTracePage />
        </MemoryRouter>
      );

      expect(
        screen.queryByRole("button", { name: UI_TEXT_MORE })
      ).not.toBeInTheDocument();
    }
  );

  it("should invoke fetchMoreDeploymentTraces when inView and hasMore are true", () => {
    mockInView = true;
    mockQueryResult.data.pages[0].hasMore = true;

    render(
      <MemoryRouter>
        <DeploymentTracePage />
      </MemoryRouter>
    );

    expect(mockFetchMoreDeploymentTraces).toHaveBeenCalledTimes(1);
  });

  it("should not invoke fetchMoreDeploymentTraces when inView is true but hasMore is false", () => {
    mockInView = true;
    mockQueryResult.data.pages[0].hasMore = false;

    render(
      <MemoryRouter>
        <DeploymentTracePage />
      </MemoryRouter>
    );

    expect(mockFetchMoreDeploymentTraces).not.toHaveBeenCalled();
  });

  it("should not invoke fetchMoreDeploymentTraces when inView is true but isFetching is true", () => {
    mockInView = true;
    mockQueryResult.data.pages[0].hasMore = true;
    mockQueryResult.isFetching = true;

    render(
      <MemoryRouter>
        <DeploymentTracePage />
      </MemoryRouter>
    );

    expect(mockFetchMoreDeploymentTraces).not.toHaveBeenCalled();
  });

  it("should disable the More button when isFetching is true", () => {
    mockQueryResult.data.pages[0].hasMore = true;
    mockQueryResult.isFetching = true;

    render(
      <MemoryRouter>
        <DeploymentTracePage />
      </MemoryRouter>
    );

    const moreButton = screen.getByRole("button", { name: UI_TEXT_MORE });
    expect(moreButton).toBeDisabled();
  });

  it("should call refetchDeploymentTraces when refresh button is clicked", () => {
    render(
      <MemoryRouter>
        <DeploymentTracePage />
      </MemoryRouter>
    );

    const refreshButton = screen.getByRole("button", { name: UI_TEXT_REFRESH });
    fireEvent.click(refreshButton);
    expect(mockRefetchDeploymentTraces).toHaveBeenCalledTimes(1);
  });

  it("should display empty state when there are no traces and not fetching", () => {
    mockQueryResult.data.pages = [];

    render(
      <MemoryRouter>
        <DeploymentTracePage />
      </MemoryRouter>
    );

    expect(screen.getByText("No deployments")).toBeInTheDocument();
  });
});
