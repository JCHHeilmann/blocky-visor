export interface TooltipData {
  html: string;
  x: number;
  y: number;
}

function createTooltipStore() {
  let current = $state<TooltipData | null>(null);

  return {
    get current() {
      return current;
    },
    show(html: string, x: number, y: number) {
      current = { html, x, y };
    },
    hide() {
      current = null;
    },
  };
}

export const tooltipStore = createTooltipStore();

export function tooltip(node: HTMLElement | SVGElement, text: () => string) {
  function onEnter(e: Event) {
    const me = e as MouseEvent;
    tooltipStore.show(text(), me.clientX, me.clientY);
  }

  function onMove(e: Event) {
    const me = e as MouseEvent;
    if (tooltipStore.current) {
      tooltipStore.show(tooltipStore.current.html, me.clientX, me.clientY);
    }
  }

  function onLeave() {
    tooltipStore.hide();
  }

  node.addEventListener("mouseenter", onEnter);
  node.addEventListener("mousemove", onMove);
  node.addEventListener("mouseleave", onLeave);

  return {
    update(newText: () => string) {
      text = newText;
    },
    destroy() {
      node.removeEventListener("mouseenter", onEnter);
      node.removeEventListener("mousemove", onMove);
      node.removeEventListener("mouseleave", onLeave);
      tooltipStore.hide();
    },
  };
}
