import CategorySummaries from './CategorySummaries';
import NavBar from './NavBar';
import {
  BrowserRouter as Router,
  Routes,
  Route
} from "react-router-dom";


function App() {
  return (
    <>
      <Router basename={'/app'}>
        <Routes>
          <Route path="/new" element={<NewTxn />}>
          </Route>
          <Route path="/txns" element={<Txns />}>
          </Route>
          <Route path="/" element={<CategorySummaries />}>
          </Route>
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

export default App;
