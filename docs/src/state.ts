import React, { useReducer, Reducer, ReducerState } from 'react';

interface State {
   [key: string]: {
    error?: any,
    loading?: boolean,
    value: any,
   }
};

interface Action {
  type: string,
  loading?: boolean,
  error?: any,
  value?: any,
};

interface ContextIface {
  state: ReducerState<Reducer<State, Action>>,
  dispatch(type: string, promise: Promise<any>): void
  setState(type: string, value: any): void
  setError(type: string, error: any): void
  setLoading(type: string, loading: boolean): void
}
let Context: React.Context<ContextIface>;

function useDispatcher() {
  let setState = (type: string, value: any) => { };
  let setError = (type: string, error: any) => { };
  let setLoading = (type: string, loading: boolean) => { };

  const [state, doDispatch] = useReducer<Reducer<State, Action>>((state: State, action: Action) => {
    const { type, error, value, loading } = action;
    let prev = state[type];
    return {
      ...state,
      [type]: {
        loading: loading ?? prev?.loading,
        error: error ?? prev?.error,
        value,
      },
    };
  }, {
  })

  setState = (type: string, value: any) => {
    doDispatch({
      type,
      value,
      loading: false,
      error: undefined,
    });
  }
  setLoading = (type: string, loading: boolean) => {
    doDispatch({
      type,
      loading
    });
  }
  setError = (type: string, error: any) => {
    doDispatch({
      type,
      error
    });
  }
  const dispatch = (type: string, promise: Promise<any>) => {
    setLoading(type, true);
    promise.then(
      result => setState(type, result)
    ).catch(
      e => setError(type, e)
    )
  }
  const result = { state, dispatch, setState, setError, setLoading }
  Context = React.createContext(result);
  return result;
};

export { Context, useDispatcher };
