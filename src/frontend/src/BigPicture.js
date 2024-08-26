import React from 'react';
import useFetch from "./useFetch";
import Header from './Header.js';
import Spinner from './Spinner.js'
import './BigPicture.css';
import Box from '@mui/material/Box';
import Typography from '@mui/material/Typography';

export default function BigPicture() {
    const { response, error, loading } = useFetch(
        "/api/bigpicture/",
        {
            query: {},
        }
    );

    if (loading) {
        return (
            <>
                <Header back_visibility="visible" title="Big Picture"></Header>
                <Spinner />
            </>
        );
    }
    if (error) {
        return <div className="error">{JSON.stringify(error)}</div>;
    }
    if (typeof response !== 'undefined' && "error" in response) {
        return <div className="error">Error: {response.error}</div>;
    }

    const net_worth = parseFloat(response.net_worth)

    const taxes_twelve_months = parseFloat(response.taxes_twelve_months)
    const taxes_three_months = parseFloat(response.taxes_three_months)

    // 'Income' means net income, after taxes; we also exclude taxes from
    // 'Expenses'.
    const income_three_months = parseFloat(response.income_three_months) + taxes_three_months
    const expenses_three_months = parseFloat(response.expenses_three_months) - taxes_three_months
    const income_twelve_months = parseFloat(response.income_twelve_months) + taxes_twelve_months
    const expenses_twelve_months = parseFloat(response.expenses_twelve_months) - taxes_twelve_months

    const current_year = new Date().getFullYear()

    const savingsNeededToRetire = (-1 * expenses_twelve_months) / 0.04
    const annualAmountSaved = (income_twelve_months + expenses_twelve_months)
    const assumedInterestRate = 0.04
    const discountRate = assumedInterestRate / (1 + assumedInterestRate)
    const yearsToSave = Math.log((annualAmountSaved + savingsNeededToRetire * discountRate) / (annualAmountSaved + discountRate * net_worth)) / Math.log(assumedInterestRate + 1)

    return (
        <>
            <Header back_visibility="visible" title="Big Picture"></Header>
            <Box sx={{ p: 2, pb: 8 }}>
                <Typography variant="body1" component="div" align="center" gutterBottom>
                    Our net worth is currently <span class="bigPictureNumber">${formatDollars(net_worth)}</span>
                </Typography>
                <br />
                <Typography variant="h6" component="div" align="center" gutterBottom>
                    Last three months
                </Typography>
                <Typography variant="body1" component="div" align="center" gutterBottom>
                    In the last three months, we've
                    <table>
                        <tr>
                            <td>made</td>
                            <td>${formatDollars(income_three_months)}</td>
                        </tr>
                        <tr>
                            <td>spent</td>
                            <td>${formatDollars(-1 * expenses_three_months)}</td>
                        </tr>
                    </table>
                </Typography>
                <Typography variant="body1" component="div" align="center" gutterBottom>
                    Savings rate: <span class="bigPictureNumber">{Math.round(((income_three_months + expenses_three_months) / income_three_months) * 100)}% of income saved</span>
                </Typography>
                <br />
                <Typography variant="h6" component="div" align="center" gutterBottom>
                    Last twelve months
                </Typography>
                <Typography variant="body1" component="div" align="center" gutterBottom>
                    In the last twelve months, we've
                    <table>
                        <tr>
                            <td>made</td>
                            <td>${formatDollars(income_twelve_months)}</td>
                        </tr>
                        <tr>
                            <td>spent</td>
                            <td>${formatDollars(-1 * expenses_twelve_months)}</td>
                        </tr>
                    </table>
                </Typography>
                <Typography variant="body1" component="div" align="center" gutterBottom>
                    Savings rate: <span class="bigPictureNumber">{Math.round(((income_twelve_months + expenses_twelve_months) / income_twelve_months) * 100)}% of income saved</span>
                </Typography>
                <br />
                <Typography variant="h6" component="div" align="center" gutterBottom>
                    Retirement
                </Typography>
                <Typography variant="body1" component="div" align="center" gutterBottom>
                    Based on our 12-month expenses, we need ${formatDollars(Math.round(savingsNeededToRetire / 1000) * 1000)} to retire.
                </Typography>
                <Typography variant="body1" component="div" align="center" gutterBottom>
                    Based on our 12-month savings rate, our savings can replace our income in <span class="bigPictureNumber">{Math.round(yearsToSave)} years</span> (in {current_year + Math.round(yearsToSave)})
                </Typography>
            </Box>
        </>
    );
}

function formatDollars(x) {
    return Math.round(x).toString().replace(/\B(?=(\d{3})+(?!\d))/g, ",");
}
