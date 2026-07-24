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
  // biome-ignore lint/complexity/noBannedTypes: Next.js built-in type signature
  type NextLayoutComponentType<P = {}> = NextComponentType<
    NextPageContext,
    any,
    P
  > & {
    getLayout?: (page: ReactNode) => ReactNode;
  };

  // biome-ignore lint/complexity/noBannedTypes: Next.js built-in type signature
  type NextLayoutPage<P = {}, IP = P> = NextComponentType<
    NextPageContext,
    IP,
    P
  > & {
    getLayout: (page: ReactNode) => ReactNode;
  };
}

declare module "next/app" {
  // biome-ignore lint/complexity/noBannedTypes: Next.js built-in type signature
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

export type IDiscordRole = {
  ID: number;
  NativeId: string;
  Name: string;
  Mentionable: number;
  Hoist: number;
  Color: number;
  ImplicitType: number;
};

export type IDiscordUserRoleAssignment = {
  ID: number;
  UserNativeID: string;
  RoleNativeID: string;
  GrantedDate: string;
  ExpirationDate: string;
  Status: "active" | "expired";
  TimeRemaining: string;
};

export type IDiscordChannel = {
  id: string;
  name: string;
  type: number;
  position: number;
};

export type IDiscordEmoji = {
  id: string;
  name: string;
  animated: boolean;
  image_url: string;
  api_name: string;
};

export type IDiscordMember = {
  native_id: string;
  username: string;
  nickname: string;
  avatar: string;
};

export type IDiscordRoleReactionEmbed = {
  ID: number;
  NativeMessageId: string;
  Name: string;
  Payload: string;
  Mode: string;
  Tags: string;
  Version: number;
  CreatedAt: string;
  UpdatedAt: string;
};

export type IReactionRoleMessagePayload = {
  channel_id: string;
  message?: string;
  mode?: "default" | "reverse";
  embed?: IReactionRoleEmbed;
  interactions: IReactionInteraction[];
};

export type IReactionRoleEmbed = {
  title?: string;
  description?: string;
  color?: number;
  image_url?: string;
  thumbnail_url?: string;
  footer?: string;
  author?: string;
  fields?: IReactionEmbedField[];
};

export type IReactionEmbedField = {
  name: string;
  value: string;
  inline?: boolean;
};

export type IReactionInteraction = {
  id: string;
  type: "emoji" | "button" | "dropdown";
  emoji?: string;
  label?: string;
  style?: "primary" | "secondary" | "success" | "danger";
  role_native_id?: string;
  placeholder?: string;
  options?: IDropdownOption[];
};

export type IDropdownOption = {
  id: string;
  label: string;
  emoji?: string;
  description?: string;
  role_native_id: string;
};

export type IAuditLogItem = {
  ID: number;
  NativeMessageId: string;
  ChannelId: string;
  GuildId: string;
  AuthorId: string;
  AuthorName: string;
  Action: "delete" | "edit";
  BeforeContent: string;
  AfterContent: string;
  Attachments: string;
  CreatedAt: string;
  UpdatedAt: string;
};

export type IAuditLogListResponse = {
  items: IAuditLogItem[];
  total: number;
  limit: number;
  offset: number;
};

declare namespace NodeJS {
  interface ProcessEnv {
    DISCORD_GUILD_ID: string;
  }
}
