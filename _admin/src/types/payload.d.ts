export interface IModuleConfigInput {
  module_id: number;
  new_config: Record<string, any>;
}

export interface IPinModuleConfig {
  threshold: number;
}
