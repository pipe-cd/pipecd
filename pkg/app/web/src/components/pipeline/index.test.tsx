import React from "react";
import { createStore, render } from "../../../test-utils";
import { updateActiveStage } from "../../modules/active-stage";
import { dummyDeployment } from "../../__fixtures__/dummy-deployment";

import { Pipeline } from "./";

it("should dispatch updateActiveState when first rendering", () => {
  const store = createStore({
    deployments: {
      entities: {
        [dummyDeployment.id]: dummyDeployment,
      },
      ids: [dummyDeployment.id],
    },
  });

  render(<Pipeline deploymentId={dummyDeployment.id} />, { store });

  expect(store.getActions()).toEqual([
    {
      type: updateActiveStage.type,
      payload: {
        name: dummyDeployment.stagesList[0].name,
        stageId: dummyDeployment.stagesList[0].id,
        deploymentId: dummyDeployment.id,
      },
    },
  ]);
});
