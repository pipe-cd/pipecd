import React from "react";
import { KubernetesStateView } from "./kubernetes-state-view";

export default {
  title: "APPLICATION/KubernetesStateView",
  component: KubernetesStateView,
};

export const overview: React.FC = () => (
  <KubernetesStateView
    resources={[
      {
        id: "8621f186-6641-4f7a-9be4-5983eb647f8d",
        ownerIdsList: ["660ecdfd-307b-4e47-becd-1fde4e0c1e7a"],
        parentIdsList: [],
        name: "demo-application-9504e8601a",
        apiVersion: "apps/v1",
        kind: "ReplicaSet",
        namespace: "default",
        healthStatus: 0,
        healthDescription: "",
        createdAt: 1592472088,
        updatedAt: 1592472088,
      },
      {
        id: "ae5d0031-1f63-4396-b929-fa9987d1e6de",
        ownerIdsList: ["660ecdfd-307b-4e47-becd-1fde4e0c1e7a"],
        parentIdsList: ["8621f186-6641-4f7a-9be4-5983eb647f8d"],
        name: "demo-application-9504e8601a-7vrdw",
        apiVersion: "v1",
        kind: "Pod",
        namespace: "default",
        healthStatus: 0,
        healthDescription: "",
        createdAt: 1592472088,
        updatedAt: 1592472088,
      },
      {
        id: "f55c7891-ba25-44bb-bca4-ffbc16b0089f",
        ownerIdsList: ["660ecdfd-307b-4e47-becd-1fde4e0c1e7a"],
        parentIdsList: ["8621f186-6641-4f7a-9be4-5983eb647f8d"],
        name: "demo-application-9504e8601a-vlgd5",
        apiVersion: "v1",
        kind: "Pod",
        namespace: "default",
        healthStatus: 0,
        healthDescription: "",
        createdAt: 1592472088,
        updatedAt: 1592472088,
      },
    ]}
  />
);
