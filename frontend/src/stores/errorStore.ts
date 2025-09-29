import { create } from 'zustand';

interface ErrorState {
  error: string | null;
  setError: (err: string) => void;
  clearError: () => void;
}

export const useErrorStore = create<ErrorState>((set) => ({
  error: null,
  setError: (err) => set({ error: err }),
  clearError: () => set({ error: null }),
}));