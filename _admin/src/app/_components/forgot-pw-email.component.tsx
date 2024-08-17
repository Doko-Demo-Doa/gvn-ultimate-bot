import { Button, Html } from "@react-email/components";
import * as React from "react";

interface Props {
  token: string;
}

const ForgotPasswordEmailComponent: React.FC<Props> = ({ token }) => {
  return (
    <Html>
      <Button
        href="https://example.com"
        style={{ background: "#000", color: "#fff", padding: "12px 20px" }}
      >
        Click me {token}
      </Button>
    </Html>
  );
};

export default ForgotPasswordEmailComponent;
