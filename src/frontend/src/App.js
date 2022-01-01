import CategorySummaries from './CategorySummaries';
import CategoryDetail from './CategoryDetail';
import TransactionList from './TransactionList';
import TransactionDetail from './TransactionDetail';
import NewTxn from './NewTxn';
import NavBar from './NavBar';
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
