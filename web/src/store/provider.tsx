"use client";

import { type ReactNode, createContext, useRef, useContext } from "react";
import { useStore } from "zustand";
import { type Store, createStore } from "@/store/state";

export type AppStoreApi = ReturnType<typeof createStore>;

export const AppStoreContext = createContext<AppStoreApi | undefined>(
  undefined,
);

export interface StoreProviderProps {
  children: ReactNode;
}

export const AppStoreProvider = ({ children }: StoreProviderProps) => {
  const storeRef = useRef<AppStoreApi>(createStore());
  return (
    <AppStoreContext.Provider value={storeRef.current}>
      {children}
    </AppStoreContext.Provider>
  );
};

export const useAppStore = <T,>(selector: (store: Store) => T): T => {
  const storeContext = useContext(AppStoreContext);

  if (!storeContext) {
    throw new Error(`useAppStore must be used within StoreProvider`);
  }

  return useStore(storeContext, selector);
};
