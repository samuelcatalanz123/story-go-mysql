import "@testing-library/jest-dom";

// Node 25 expone un `localStorage` experimental global que tapa al de jsdom y
// no funciona sin `--localstorage-file`. Lo sustituimos por uno en memoria para
// que los tests usen un almacenamiento real.
if (typeof localStorage === "undefined" || typeof localStorage.setItem !== "function") {
  class MemoryStorage implements Storage {
    private store = new Map<string, string>();
    get length(): number {
      return this.store.size;
    }
    clear(): void {
      this.store.clear();
    }
    getItem(key: string): string | null {
      return this.store.has(key) ? this.store.get(key)! : null;
    }
    key(index: number): string | null {
      return Array.from(this.store.keys())[index] ?? null;
    }
    removeItem(key: string): void {
      this.store.delete(key);
    }
    setItem(key: string, value: string): void {
      this.store.set(key, String(value));
    }
  }
  Object.defineProperty(globalThis, "localStorage", {
    value: new MemoryStorage(),
    configurable: true,
    writable: true,
  });
}
