import { useState, useEffect } from "preact/hooks";
import { dateRepr, dateTimeString, relRepaintDelay } from "../utils";
import Tooltip from "./Tooltip";

interface RelativeTimeProps {
  val: string;
  locale: string;
}

export default function RelativeTime({ val, locale }: RelativeTimeProps) {
  const [formatted, setFormatted] = useState("");
  const date = new Date(val);
  const title = dateTimeString(date, locale);

  useEffect(() => {
    let timer: number | undefined;

    function repaint() {
      setFormatted(dateRepr(date, locale));
      const delay = relRepaintDelay(date);
      if (delay !== null) {
        timer = window.setTimeout(repaint, delay);
      }
    }

    repaint();

    return () => {
      if (timer !== undefined) window.clearTimeout(timer);
    };
  }, [val, locale]);

  return (
    <Tooltip label={title}>
      <time dateTime={val}>{formatted}</time>
    </Tooltip>
  );
}
