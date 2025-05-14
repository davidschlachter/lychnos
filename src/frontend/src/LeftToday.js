import { useTheme } from '@mui/material/styles';

const red = '#ff5e00';
const green = '#a2ff00';

function LeftToday(props) {
    const theme = useTheme();

    let sign = '';
    let color = theme.palette.text.primary;
    if (Math.abs(props.actual) > Math.abs(props.budgetted * props.timeSpent / 100)) {
        sign = '-';
        if (props.budgetted > 0) {
            color = green; // Celebrate income categories greater than target
        } else {
            color = red;
        }
    }
    const underOver = Math.round(Math.abs(props.actual - (props.budgetted * props.timeSpent / 100)));

    return (
        <span style={{
            "color": color
        }}>{sign}{underOver}</span>
    );
}

export default LeftToday;
