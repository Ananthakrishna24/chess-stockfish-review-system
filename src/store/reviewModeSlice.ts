import { createSlice, PayloadAction } from '@reduxjs/toolkit';

interface ReviewModeState {
  isReviewMode: boolean;
  showMoveClassifications: boolean;
  autoPlayMode: boolean;
  autoPlayInterval: number; // milliseconds
  showAlternativeMoves: boolean;
  highlightMistakes: boolean;
  currentReviewMove: number;
}

const initialState: ReviewModeState = {
  isReviewMode: false,
  showMoveClassifications: true,
  autoPlayMode: false,
  autoPlayInterval: 2000,
  showAlternativeMoves: true,
  highlightMistakes: true,
  currentReviewMove: 0,
};

const reviewModeSlice = createSlice({
  name: 'reviewMode',
  initialState,
  reducers: {
    startReviewMode: (state) => {
      state.isReviewMode = true;
      state.currentReviewMove = 0;
    },
    exitReviewMode: (state) => {
      state.isReviewMode = false;
      state.autoPlayMode = false;
      state.currentReviewMove = 0;
    },
    toggleAutoPlay: (state) => {
      state.autoPlayMode = !state.autoPlayMode;
    },
    setAutoPlayInterval: (state, action: PayloadAction<number>) => {
      state.autoPlayInterval = action.payload;
    },
    toggleMoveClassifications: (state) => {
      state.showMoveClassifications = !state.showMoveClassifications;
    },
    toggleAlternativeMoves: (state) => {
      state.showAlternativeMoves = !state.showAlternativeMoves;
    },
    toggleHighlightMistakes: (state) => {
      state.highlightMistakes = !state.highlightMistakes;
    },
    setCurrentReviewMove: (state, action: PayloadAction<number>) => {
      state.currentReviewMove = action.payload;
    },
  },
});

export const {
  startReviewMode,
  exitReviewMode,
  toggleAutoPlay,
  setAutoPlayInterval,
  toggleMoveClassifications,
  toggleAlternativeMoves,
  toggleHighlightMistakes,
  setCurrentReviewMove,
} = reviewModeSlice.actions;

export default reviewModeSlice.reducer; 