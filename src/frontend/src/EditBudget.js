import React from 'react';
import EditBudgetForm from './EditBudgetForm.js';
import Header from './Header.js';
import Spinner from './Spinner.js'
import axios from 'axios';

export default function EditBudget() {
    const [resp, setResp] = React.useState({ categorybudgets: null, categories: null, seedBudget: null, loading: true });
    React.useEffect(() => {
        const fetchData = async () => {
            const catsGlobal = await axios("/api/categories/");
            const cbsGlobal = await axios("/api/categorybudgets/");
            const seed = (function () {
                let bgt = {};
                for (const c of catsGlobal.data) {
                    bgt[c.id] = getAmount(c.id, cbsGlobal.data);
                }
                return bgt;
            })();
            setResp({ categorybudgets: cbsGlobal.data, categories: catsGlobal.data, seedBudget: seed, loading: false });
        };
        fetchData();
    }, []);

    if (resp.loading === true) {
        return (
            <>
                <Header back_visibility="visible" back_location="/" title="Edit budget"></Header>
                <Spinner />
            </>
        );
    } else {
        return (
            <>
                <Header back_visibility="visible" back_location="/" title="Edit budget" budgetedit={false}></Header>
                <EditBudgetForm categorybudgets={resp.categorybudgets} categories={resp.categories} seed={resp.seedBudget} />
            </>
        );
    }

}

function getAmount(categoryID, categorybudgets) {
    for (const cb of categorybudgets) {
        if (cb.category === categoryID) {
            return cb.amount;
        }
    }
    return 0;
}