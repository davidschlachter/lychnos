import CategoryDetail from './CategoryDetail';
import CategorySummaries from './CategorySummaries';
import EditBudget from './EditBudget';
import NavBar from './NavBar';
import NewTxn from './NewTxn';
import TransactionDetail from './TransactionDetail';
import TransactionList from './TransactionList';
import * as React from 'react';
import useMediaQuery from '@mui/material/useMediaQuery';
import { createTheme, ThemeProvider } from '@mui/material/styles';
import CssBaseline from '@mui/material/CssBaseline';
import {
  BrowserRouter as Router,
  Routes,
  Route,
  useParams
} from "react-router-dom";


function App() {

  const prefersDarkMode = useMediaQuery('(prefers-color-scheme: dark)');

  const theme = React.useMemo(
    () =>
      createTheme({
        palette: {
          mode: prefersDarkMode ? 'dark' : 'light',
        },
      }),
    [prefersDarkMode],
  );

  return (
    <ThemeProvider theme={theme}>
      <CssBaseline />
      <Router basename={'/app'}>
        <Routes>
          <Route path="/newTxn" element={<NewTxn />} />
          <Route path="/txns/:categoryID" element={<TransactionList />} />
          <Route path="/txns" element={<TransactionList />} />
          <Route path="/" element={<CategorySummaries />} />
          <Route path="/categorydetail/:categoryId" element={<CategoryDetailHelper />} />
          <Route path="/txn/:txnID" element={<TransactionDetail />} />
          <Route path="/budget/" element={<EditBudget />} />
        </Routes>
        <NavBar />
      </Router>
    </ThemeProvider>
  );
}

function CategoryDetailHelper() {
  const { categoryId } = useParams();
  return <CategoryDetail categoryId={categoryId} />
}

export default App;
