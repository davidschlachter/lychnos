import CategorySummaries from './CategorySummaries';
import CategoryDetail from './CategoryDetail';
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
          <Route path="/new" element={<NewTxn />} />
          <Route path="/txns" element={<Txns />} />
          <Route path="/" element={<CategorySummaries />} />
          <Route path="/categorydetail/:categoryId" element={<CategoryDetailHelper />} />
        </Routes>
        <NavBar />
      </Router>
    </>
  );
}

function NewTxn() {
  return <h2>New transaction</h2>;
}

function Txns() {
  return <h2>List of the transactions</h2>;
}

function CategoryDetailHelper() {
  const { categoryId } = useParams();
  return <CategoryDetail categoryId={categoryId} />
}

export default App;
