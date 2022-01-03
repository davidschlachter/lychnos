import CategoryDetail from './CategoryDetail';
import CategorySummaries from './CategorySummaries';
import EditBudget from './EditBudget';
import NavBar from './NavBar';
import NewTxn from './NewTxn';
import TransactionDetail from './TransactionDetail';
import TransactionList from './TransactionList';
import {
  BrowserRouter as Router,
  Routes,
  Route,
  useParams
} from "react-router-dom";


function App() {
  return (
    <>
      <Router basename={'/app'}>
        <Routes>
          <Route path="/newTxn" element={<NewTxn />} />
          <Route path="/txns" element={<TransactionList />} />
          <Route path="/" element={<CategorySummaries />} />
          <Route path="/categorydetail/:categoryId" element={<CategoryDetailHelper />} />
          <Route path="/txn/:txnID" element={<TransactionDetail />} />
          <Route path="/budget/" element={<EditBudget />} />
        </Routes>
        <NavBar />
      </Router>
    </>
  );
}

function CategoryDetailHelper() {
  const { categoryId } = useParams();
  return <CategoryDetail categoryId={categoryId} />
}

export default App;
