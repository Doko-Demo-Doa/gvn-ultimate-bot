declare global {
  namespace NodeJS {
    interface ProcessEnv {
      BACKEND_BASE_URL: string;
    }
  }

  enum ModuleStatus {
    DISABLED,
    ENABLED,
  }

  type ModuleActivationStatusType = {
    [moduleName: string]: ModuleStatus;
  };
}

export {};
