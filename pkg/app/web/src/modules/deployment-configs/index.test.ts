import {
  deploymentConfigsSlice,
  DeploymentConfigsState,
  clearTemplateTarget,
  fetchTemplateList,
  selectTemplatesByAppId,
} from "./";
import { dummyDeploymentConfigTemplates } from "~/__fixtures__/dummy-deployment-config";
import { addApplication } from "../applications";

const initialState: DeploymentConfigsState = {
  templates: {},
  targetApplicationId: null,
};

describe("deploymentConfigsSlice reducer", () => {
  it("should return the initial state", () => {
    expect(
      deploymentConfigsSlice.reducer(undefined, {
        type: "TEST_ACTION",
      })
    ).toEqual(initialState);
  });

  it(`should handle ${clearTemplateTarget.type}`, () => {
    expect(
      deploymentConfigsSlice.reducer(
        { ...initialState, targetApplicationId: "application-1" },
        {
          type: clearTemplateTarget.type,
        }
      )
    ).toEqual(initialState);
  });

  describe("fetchTemplateList", () => {
    it(`should handle ${fetchTemplateList.fulfilled.type}`, () => {
      expect(
        deploymentConfigsSlice.reducer(
          { templates: {}, targetApplicationId: "application-1" },
          {
            type: fetchTemplateList.fulfilled.type,
            payload: dummyDeploymentConfigTemplates,
          }
        )
      ).toEqual({
        targetApplicationId: "application-1",
        templates: {
          "application-1": dummyDeploymentConfigTemplates,
        },
      });
    });
  });

  describe("addApplication", () => {
    it(`should handle ${addApplication.fulfilled.type}`, () => {
      expect(
        deploymentConfigsSlice.reducer(initialState, {
          type: addApplication.fulfilled.type,
          payload: "application-id",
        })
      ).toEqual({
        ...initialState,
        targetApplicationId: "application-id",
      });
    });
  });
});

describe("selectTemplatesByAppId", () => {
  let state: DeploymentConfigsState;
  beforeEach(() => {
    state = {
      targetApplicationId: null,
      templates: { "app-1": dummyDeploymentConfigTemplates },
    };
  });

  it("should return null if target is null", () => {
    expect(selectTemplatesByAppId(state)).toEqual(null);
  });

  it("should return templates by target id", () => {
    state.targetApplicationId = "app-1";
    expect(selectTemplatesByAppId(state)).toEqual(
      dummyDeploymentConfigTemplates
    );
  });
});
