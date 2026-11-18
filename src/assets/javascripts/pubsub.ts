type Event = {
  selectFeed: [feed: number];
};

type Handler<T extends unknown[]> = (...args: T) => void;

class PubSub {
  private handlers = new Map<keyof Event, Set<Handler<any[]>>>();

  sub<K extends keyof Event>(event: K, handler: Handler<Event[K]>): () => void {
    let set = this.handlers.get(event);
    if (!set) {
      set = new Set();
      this.handlers.set(event, set);
    }
    set.add(handler);
    return () => {
      set.delete(handler);
    };
  }

  pub<K extends keyof Event>(event: K, ...args: Event[K]): void {
    this.handlers.get(event)?.forEach(handler => handler(...args));
  }
}

export const events = new PubSub();
