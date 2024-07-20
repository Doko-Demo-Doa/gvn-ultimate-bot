"use client";

import { useAppModules } from "~/hooks/use-app-modules";

const moduleNameMap: Record<string, string> = {
  pin_module: "Pin Module",
  grant_role_module: "Grant Role Module",
};

const EnabledModules = () => {
  const { data } = useAppModules();

  if (!data) {
    return <div>Loading...</div>;
  }

  return (
    <div>
      <h1>Enabled Modules</h1>

      {data?.data.map((module) => (
        <div key={module.ID}>
          <h2>{moduleNameMap[module.ModuleName] || ""}</h2>
        </div>
      ))}
    </div>
  );
};

export default EnabledModules;
