// NextJS typing extensions
import type {
  NextComponentType,
  NextLayoutComponentType,
  NextPageContext,
} from "next";
import type { AppProps } from "next/app";

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

declare module "next" {
  type NextLayoutComponentType<P = {}> = NextComponentType<
    NextPageContext,
    any,
    P
  > & {
    getLayout?: (page: ReactNode) => ReactNode;
  };

  type NextLayoutPage<P = {}, IP = P> = NextComponentType<
    NextPageContext,
    IP,
    P
  > & {
    getLayout: (page: ReactNode) => ReactNode;
  };
}

declare module "next/app" {
  type AppLayoutProps<P = {}> = AppProps & {
    Component: NextLayoutComponentType;
  };
}

declare module "*.json" {
  const value: any;
  export default value;
}

declare module "iron-session" {
  interface IronSessionData {
    nonce?: string;
    siwe?: SiweMessage;
  }
}

enum ModuleActivation {
  DISABLED = 0,
  ENABLED = 1,
}

export type BackendResponseType<D> = {
  code: number;
  message: string;
  data: D;
};

export type IBackendModuleType = {
  ID: number;
  ModuleName: string;
  ModuleLabel: string;
  IsActivated: ModuleActivation;
  CustomConfig: string;
};

export interface IReactionRoleMessagePayload {
  guildId: string;
  channelId: string;
  messageId: string;
  detail: {
    message: string;
    embedTitle: string;
    embedDescription: string;
  };
}

declare namespace NodeJS {
  interface ProcessEnv {
    DISCORD_GUILD_ID: string;
  }
}
