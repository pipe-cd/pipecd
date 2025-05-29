- Start Date: 2025-04-27
- Target Version: 1.0

# Summary

This RFC proposes replacing Redux with `react-query` and `react context` to simplify the codebase, improve developer experience, and reduce boilerplate.

# Motivation

The current use of Redux introduces significant boilerplate and complexity, which can be simplified by adopting modern alternatives like React Context, hooks, and react-query. The key motivations for this change are:

- **Reduced Boilerplate**: Redux requires actions, reducers, and middleware, which can be replaced with simpler patterns using React Context and hooks.
- **Improved Developer Experience**: React Context and hooks provide a more intuitive API, making state management easier to understand and maintain.
- **Optimized Data Fetching**: react-query offers powerful features like caching, background updates, and automatic retries, which are not natively available in Redux.
- **Modern Best Practices**: React Context and hooks align with modern React development practices, reducing the need for external dependencies.
- **Performance Improvements**: By leveraging react-query for data fetching and caching, we can reduce unnecessary re-renders and improve application performance.
- **Simplified Codebase**: Removing Redux will reduce the overall complexity of the codebase, making it easier for new developers to onboard and contribute.

# Detailed design

## What should be removed:

- **Thunk Middleware**: Remove the middleware used for handling side effects in Redux.
- **Redux Context**: Eliminate the Redux provider and related context setup.
- **Module Folder**: Remove the folder containing Redux logic, including actions, reducers, and selectors.

## What should be added:

- **React Context**: Create custom contexts and providers for managing authentication state.
- **React Query**: Use `react-query` for handling API calls and caching authentication-related data.

## New Web Folder Structure

The updated folder structure will be as follows:

```
src/
├── api/          # Define all API service functions
├── queries/      # Add custom hooks for querying and mutating resources using `react-query`
├── contexts/     # Implement custom React contexts and providers for managing state
```

### Current Resources

With the new folder structure, the codebase will be organized around specific resources. The following resources are currently available in the project:

- **api-keys**
- **applications**
- **client**
- **commands**
- **deployments**
- **deploymentTraces**
- **events**
- **insight**
- **me**
- **piped**
- **project**
- **stage-log**

## Testing

Changing Redux to the new state management solution will require updating numerous tests. The following points should be considered for testing:

### Specific Areas tests need to be updated

- **Queries**:
  - Add tests for custom hooks in the `queries` folder to validate data fetching and caching behavior.
- **Components**:
  - Update tests for components that consume data from React Context or `react-query`.
- **Pages**:
  - Verify that pages render correctly and handle state changes properly.
- **Utility Functions**:
  - Update tests for utility functions that interact with Redux or API calls.

### Test Coverage

- Ensure that test coverage remains high by adding tests for all new features and changes.

# Unresolved questions

This design is complete, and no unresolved questions remain.
