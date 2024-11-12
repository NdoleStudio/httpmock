import { createStore as createZustandStore } from "zustand/vanilla";
import {
  EntitiesProject,
  RequestsProjectCreateRequest,
  ResponsesOkEntitiesProject,
} from "@/api/model";
import axios from "@/api/axios";
import { AxiosError } from "axios";
import { getErrorMessages } from "@/utils/errors";

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
    createProject: async (request) => {
      return new Promise<EntitiesProject>((resolve, reject) => {
        axios
          .post<ResponsesOkEntitiesProject>(`/v1/projects/`, request)
          .then((response) => {
            resolve(response.data.data);
          })
          .catch(async (error: AxiosError) => {
            reject(getErrorMessages(error));
          });
      });
    },
  }));
};
