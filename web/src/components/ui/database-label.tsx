import {
  SINGLESTORE_PURPLE_500,
  SINGLESTORE_PURPLE_700,
} from "@/consts/config";
import { SnowflakeSmallLogo } from "../logo/snowflake-small";
import { SingleStoreSmallLogo } from "../logo/singlestore-small";
import { SingleStoreLogo } from "../logo/singlestore";

export interface DatabaseResultLabelProps {
  database: string;
  latency: number;
}

export function DatabaseResultLabel({
  database,
  latency,
}: DatabaseResultLabelProps) {
  const getLatencyString = (latency: number) => {
    if (latency < 1000) {
      return `${latency}ms`;
    } else if (latency < 60000) {
      return `${(latency / 1000).toFixed(2)}s`;
    } else {
      const minutes = Math.floor(latency / 60000);
      const seconds = ((latency % 60000) / 1000).toFixed(2);
      return `${minutes}m ${seconds}s`;
    }
  };

  const getDatabaseString = (database: string) => {
    if (database === "snowflake") {
      return "Snowflake";
    } else if (database === "singlestore") {
      return "SingleStore";
    }
  };

  const getDatabaseLogo = (database: string) => {
    if (database === "snowflake") {
      return <SnowflakeSmallLogo size={18} />;
    } else if (database === "singlestore") {
      return <SingleStoreSmallLogo size={18} />;
    }
  };

  const getDatabaseColor = (database: string) => {
    if (database === "snowflake") {
      return "#29B5E8";
    } else if (database === "singlestore") {
      return SINGLESTORE_PURPLE_500;
    }
  };

  const LoadingIndicator = () => {
    return (
      <div className="mx-1 inline-flex animate-spin items-center text-gray-400">
        {getDatabaseLogo(database)}
      </div>
    );
  };

  if (latency === 0) {
    return LoadingIndicator();
  }

  return (
    <p className="flex flex-row items-center text-sm text-gray-400">
      <div
        style={{ color: getDatabaseColor(database) }}
        className="mx-1 inline-flex items-center"
      >
        {getDatabaseLogo(database)}
      </div>
      {getLatencyString(latency)}
    </p>
  );
}
