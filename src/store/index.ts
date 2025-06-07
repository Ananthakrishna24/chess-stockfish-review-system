import { configureStore } from '@reduxjs/toolkit';
import reviewModeSlice from './reviewModeSlice';

export const store = configureStore({
  reducer: {
    reviewMode: reviewModeSlice,
  },
});

export type RootState = ReturnType<typeof store.getState>;
export type AppDispatch = typeof store.dispatch; 