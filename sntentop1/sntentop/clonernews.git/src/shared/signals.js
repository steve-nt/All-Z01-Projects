// Signals provide the minimal reactive primitive used by later features without introducing a state library.
let activeEffect = null;

// Effect cleanup removes stale subscriptions so reactive updates do not accumulate memory leaks.
const cleanupEffect = (effect) => {
  for (const dependencies of effect.dependencies) {
    dependencies.delete(effect);
  }

  effect.dependencies.clear();
};

// Signals expose a tiny getter/setter API so dependency tracking stays explicit and predictable.
export const createSignal = (initialValue) => {
  let value = initialValue;
  const subscribers = new Set();

  // Reads register the current effect so future writes can re-run only the dependent computations.
  const get = () => {
    if (activeEffect) {
      subscribers.add(activeEffect);
      activeEffect.dependencies.add(subscribers);
    }

    return value;
  };

  // Writes fan out only when the value truly changed, which avoids unnecessary re-renders.
  const set = (nextValue) => {
    const resolvedValue = typeof nextValue === 'function' ? nextValue(value) : nextValue;

    if (Object.is(resolvedValue, value)) {
      return value;
    }

    value = resolvedValue;

    for (const subscriber of [...subscribers]) {
      if (typeof subscriber === 'function') {
        subscriber(value);
        continue;
      }

      subscriber.run();
    }

    return value;
  };

  return {
    get,
    set,
    subscribe(callback) {
      subscribers.add(callback);

      return () => subscribers.delete(callback);
    },
  };
};

// Effects wrap side effects in dependency tracking so callers can cleanly tear them down later.
export const createEffect = (effectFn) => {
  const effect = {
    dependencies: new Set(),
    cleanup: undefined,
    run() {
      cleanupEffect(effect);

      if (typeof effect.cleanup === 'function') {
        effect.cleanup();
      }

      const previousEffect = activeEffect;
      activeEffect = effect;
      effect.cleanup = effectFn();
      activeEffect = previousEffect;
    },
  };

  effect.run();

  return () => {
    cleanupEffect(effect);

    if (typeof effect.cleanup === 'function') {
      effect.cleanup();
    }
  };
};

// Computed signals are derived values that stay synchronized with their source signals automatically.
export const createComputed = (computeFn) => {
  const computed = createSignal(computeFn());

  createEffect(() => {
    computed.set(computeFn());
  });

  return computed;
};
