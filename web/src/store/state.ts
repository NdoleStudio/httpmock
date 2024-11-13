"use client";

import { createStore as createZustandStore } from "zustand/vanilla";
import {
  EntitiesProject,
  RequestsProjectCreateRequest,
  ResponsesOkEntitiesProject,
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
  createProject: (
    request: RequestsProjectCreateRequest,
  ) => Promise<EntitiesProject>;
};

export type Store = State & Actions;

export const defaultInitState: State = {
  notificationMessage: "",
};

export const createStore = (initState: State = defaultInitState) => {
  return createZustandStore<Store>()((set, get) => ({
    ...initState,
    createProject: (request): Promise<EntitiesProject> => {
      return new Promise<EntitiesProject>((resolve, reject) => {
        console.log("createProject", request);
        axios
          .post<ResponsesOkEntitiesProject>(`/v1/projects`, request)
          .then((response) => {
            resolve(response.data.data);
          })
          .catch(async (error: AxiosError<ResponsesUnprocessableEntity>) => {
            toast.error(
              error.response?.data.message ??
                "Error while creating your mock project",
            );
            reject(getErrorMessages(error));
          });
      });
    },
  }));
};
