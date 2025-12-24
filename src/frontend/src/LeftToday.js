import './LeftToday.css';

function LeftToday(props) {
    let sign = '';
    let budgetStatus = 'budgetStatusOkay';

    const actualAbs = Math.abs(props.actual)
    const budgettedAbs = Math.abs(props.budgetted)
    const timeSpentFraction = (props.timeSpent / 100)
    const timeSpentPlusOneMonthFraction = timeSpentFraction + (1 / 12)

    // If the category is beyond where we would expect it to be right now...
    if (actualAbs > budgettedAbs * timeSpentFraction) {
        sign = '−';
        budgetStatus = 'budgetStatusCaution';
        if (props.budgetted > 0) {
            // Celebrate income categories greater than target
            budgetStatus = 'budgetStatusExceedExpectations';
        } else if (actualAbs > budgettedAbs * timeSpentPlusOneMonthFraction) {
            // Only make the text red if we're outside four weeks of spending.
            // This will allow us to pay a monthly expense like rent without the
            // category going red every time.
            budgetStatus = 'budgetStatusDanger';
        }
    }
    const underOver = Math.round(
        Math.abs(
            props.actual - (props.budgetted * timeSpentFraction)
        )
    );

    return (
        <span className={budgetStatus}>{sign}{underOver}</span>
    );
}

export default LeftToday;
