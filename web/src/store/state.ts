"use client";

import { createStore as createZustandStore } from "zustand/vanilla";
import {
  EntitiesProject,
  EntitiesProjectEndpoint,
  RequestsProjectCreateRequest,
  RequestsProjectEndpointCreateRequest,
  RequestsProjectEndpointUpdateRequest,
  RequestsProjectUpdateRequest,
  ResponsesNoContent,
  ResponsesOkArrayEntitiesProject,
  ResponsesOkArrayEntitiesProjectEndpoint,
  ResponsesOkEntitiesProject,
  ResponsesOkEntitiesProjectEndpoint,
  ResponsesUnprocessableEntity,
} from "@/api/model";
import axios from "@/api/axios";
import { AxiosError } from "axios";
import { getErrorMessages } from "@/utils/errors";
import { toast } from "sonner";

export type State = {
  notificationMessage?: string;
};

export type Actions = {
  storeProject: (
    request: RequestsProjectCreateRequest,
  ) => Promise<EntitiesProject>;
  updateProject: (
    projectId: string,
    request: RequestsProjectCreateRequest,
  ) => Promise<EntitiesProject>;
  showProject: (projectId: string) => Promise<EntitiesProject>;
  deleteProject: (projectId: string) => Promise<void>;
  indexProjects: () => Promise<Array<EntitiesProject>>;
  storeProjectEndpoint: (
    projectId: string,
    request: RequestsProjectEndpointCreateRequest,
  ) => Promise<EntitiesProjectEndpoint>;
  updateProjectEndpoint: (
    projectId: string,
    projectEndpointId: string,
    request: RequestsProjectEndpointUpdateRequest,
  ) => Promise<EntitiesProjectEndpoint>;
  indexProjectEndpoint: (
    projectId: string,
  ) => Promise<Array<EntitiesProjectEndpoint>>;
  showProjectEndpoint: (
    projectId: string,
    endpointId: string,
  ) => Promise<EntitiesProjectEndpoint>;
  deleteProjectEndpoint: (
    projectId: string,
    projectEndpointId: string,
  ) => Promise<void>;
};

export type Store = State & Actions;

export const defaultInitState: State = {
  notificationMessage: "",
};

export const createStore = (initState: State = defaultInitState) => {
  return createZustandStore<Store>()(() => ({
    ...initState,
    storeProject: (
      request: RequestsProjectCreateRequest,
    ): Promise<EntitiesProject> => {
      return new Promise<EntitiesProject>((resolve, reject) => {
        axios
          .post<ResponsesOkEntitiesProject>(`/v1/projects`, request)
          .then((response) => {
            toast.success("Project created successfully.");
            resolve(response.data.data);
          })
          .catch(async (error: AxiosError<ResponsesUnprocessableEntity>) => {
            toast.error(
              error.response?.data.message ??
                "Error while creating a new project",
            );
            reject(getErrorMessages(error));
          });
      });
    },
    updateProject: (
      projectId: string,
      request: RequestsProjectUpdateRequest,
    ): Promise<EntitiesProject> => {
      return new Promise<EntitiesProject>((resolve, reject) => {
        axios
          .put<ResponsesOkEntitiesProject>(`/v1/projects/${projectId}`, request)
          .then((response) => {
            toast.success("Project updated successfully.");
            resolve(response.data.data);
          })
          .catch(async (error: AxiosError<ResponsesUnprocessableEntity>) => {
            toast.error(
              error.response?.data.message ??
                "Error while updating your project",
            );
            reject(getErrorMessages(error));
          });
      });
    },
    showProject: (projectId: string): Promise<EntitiesProject> => {
      return new Promise<EntitiesProject>((resolve, reject) => {
        axios
          .get<ResponsesOkEntitiesProject>(`/v1/projects/${projectId}`)
          .then((response) => {
            resolve(response.data.data);
          })
          .catch(async (error: AxiosError<ResponsesUnprocessableEntity>) => {
            toast.error(
              error.response?.data.message ?? "Error while fetching project",
            );
            reject(getErrorMessages(error));
          });
      });
    },
    deleteProject: (projectId: string): Promise<void> => {
      return new Promise<void>((resolve, reject) => {
        axios
          .delete<ResponsesNoContent>(`/v1/projects/${projectId}`)
          .then(() => {
            resolve();
          })
          .catch(async (error: AxiosError<ResponsesUnprocessableEntity>) => {
            toast.error(
              error.response?.data.message ??
                "Error while deleting your project",
            );
            reject(getErrorMessages(error));
          });
      });
    },
    indexProjects: (): Promise<Array<EntitiesProject>> => {
      return new Promise<Array<EntitiesProject>>((resolve, reject) => {
        axios
          .get<ResponsesOkArrayEntitiesProject>(`/v1/projects`)
          .then((response) => {
            resolve(response.data.data);
          })
          .catch(async (error: AxiosError<ResponsesUnprocessableEntity>) => {
            toast.error(
              error.response?.data.message ?? "Error while loading projects",
            );
            reject(getErrorMessages(error));
          });
      });
    },
    storeProjectEndpoint: (
      projectId: string,
      request: RequestsProjectEndpointCreateRequest,
    ): Promise<EntitiesProjectEndpoint> => {
      return new Promise<EntitiesProjectEndpoint>((resolve, reject) => {
        axios
          .post<ResponsesOkEntitiesProjectEndpoint>(
            `/v1/projects/${projectId}/endpoints`,
            request,
          )
          .then((response) => {
            toast.success("Project endpoint created successfully.");
            resolve(response.data.data);
          })
          .catch(async (error: AxiosError<ResponsesUnprocessableEntity>) => {
            toast.error(
              error.response?.data.message ??
                "Error while creating a new project endpoint",
            );
            reject(getErrorMessages(error));
          });
      });
    },
    updateProjectEndpoint: (
      projectId: string,
      projectEndpointId: string,
      request: RequestsProjectEndpointUpdateRequest,
    ): Promise<EntitiesProjectEndpoint> => {
      return new Promise<EntitiesProjectEndpoint>((resolve, reject) => {
        axios
          .put<ResponsesOkEntitiesProjectEndpoint>(
            `/v1/projects/${projectId}/endpoints/${projectEndpointId}`,
            request,
          )
          .then((response) => {
            toast.success("Endpoint updated successfully.");
            resolve(response.data.data);
          })
          .catch(async (error: AxiosError<ResponsesUnprocessableEntity>) => {
            toast.error(
              error.response?.data.message ??
                "Error while updating a project endpoint",
            );
            reject(getErrorMessages(error));
          });
      });
    },
    indexProjectEndpoint: (
      projectId: string,
    ): Promise<Array<EntitiesProjectEndpoint>> => {
      return new Promise<Array<EntitiesProjectEndpoint>>((resolve, reject) => {
        axios
          .get<ResponsesOkArrayEntitiesProjectEndpoint>(
            `/v1/projects/${projectId}/endpoints`,
          )
          .then((response) => {
            resolve(response.data.data);
          })
          .catch(async (error: AxiosError<ResponsesUnprocessableEntity>) => {
            toast.error(
              error.response?.data.message ??
                "Error while loading project endpoints",
            );
            reject(getErrorMessages(error));
          });
      });
    },
    showProjectEndpoint: (
      projectId: string,
      projectEndpointId,
    ): Promise<EntitiesProjectEndpoint> => {
      return new Promise<EntitiesProjectEndpoint>((resolve, reject) => {
        axios
          .get<ResponsesOkEntitiesProjectEndpoint>(
            `/v1/projects/${projectId}/endpoints/${projectEndpointId}`,
          )
          .then((response) => {
            resolve(response.data.data);
          })
          .catch(async (error: AxiosError<ResponsesUnprocessableEntity>) => {
            toast.error(
              error.response?.data.message ??
                "Error while fetching project endpoint",
            );
            reject(getErrorMessages(error));
          });
      });
    },
    deleteProjectEndpoint: (
      projectId: string,
      projectEndpointId: string,
    ): Promise<void> => {
      return new Promise<void>((resolve, reject) => {
        axios
          .delete<ResponsesNoContent>(
            `/v1/projects/${projectId}/endpoints/${projectEndpointId}`,
          )
          .then(() => {
            resolve();
          })
          .catch(async (error: AxiosError<ResponsesUnprocessableEntity>) => {
            toast.error(
              error.response?.data.message ??
                "Error while deleting your endpoint",
            );
            reject(getErrorMessages(error));
          });
      });
    },
  }));
};
