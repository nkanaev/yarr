import { useSignal } from "@preact/signals";
import { useEffect, useLayoutEffect, useRef } from "preact/hooks";
import { createPortal } from "preact";
import type { ComponentChildren } from "preact";

const SHOW_DELAY = 800;
const HIDE_DELAY = 100;
const VIEWPORT_PADDING = 4;

interface TooltipProps {
  label: string;
  className?: string;
  children: ComponentChildren;
}

export default function Tooltip({ label, className = "", children }: TooltipProps) {
  const visible = useSignal(false);
  const top = useSignal(0);
  const left = useSignal(0);
  const anchor = useRef<HTMLSpanElement>(null);
  const bubble = useRef<HTMLSpanElement>(null);
  const timer = useRef<number | undefined>(undefined);

  const place = () => {
    const rect = anchor.current?.getBoundingClientRect();
    if (!rect) return;
    top.value = rect.bottom;
    left.value = rect.left + rect.width / 2;
  };

  const show = () => {
    clearTimeout(timer.current);
    timer.current = window.setTimeout(() => {
      place();
      visible.value = true;
    }, SHOW_DELAY);
  };

  const hide = () => {
    clearTimeout(timer.current);
    timer.current = window.setTimeout(() => {
      visible.value = false;
    }, HIDE_DELAY);
  };

  useLayoutEffect(() => {
    if (!visible.value || !bubble.current) return;
    const rect = bubble.current.getBoundingClientRect();
    let adjust = 0;
    if (rect.right > window.innerWidth - VIEWPORT_PADDING) {
      adjust = window.innerWidth - VIEWPORT_PADDING - rect.right;
    } else if (rect.left < VIEWPORT_PADDING) {
      adjust = VIEWPORT_PADDING - rect.left;
    }
    if (Math.abs(adjust) > 0.5) left.value += adjust;
  }, [visible.value, top.value, left.value]);

  useEffect(() => {
    if (!visible.value) return;
    const reposition = () => place();
    window.addEventListener("scroll", reposition, true);
    window.addEventListener("resize", reposition);
    return () => {
      window.removeEventListener("scroll", reposition, true);
      window.removeEventListener("resize", reposition);
    };
  }, [visible.value]);

  useEffect(() => () => clearTimeout(timer.current), []);

  return (
    <>
      <span
        ref={anchor}
        className={`c-tooltip ${className}`.trim()}
        onMouseEnter={show}
        onMouseLeave={hide}
        onFocusIn={show}
        onFocusOut={hide}>
        {children}
      </span>
      {visible.value &&
        createPortal(
          <span
            ref={bubble}
            role="tooltip"
            className="c-tooltip-bubble"
            style={{ top: `${top.value}px`, left: `${left.value}px` }}>
            {label}
          </span>,
          document.body,
        )}
    </>
  );
}
