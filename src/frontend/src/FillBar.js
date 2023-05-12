import './FillBar.css';
import useMediaQuery from '@mui/material/useMediaQuery';

const green = '#a2ff00';
const yellow = '#ffdd00';
const red = '#ff5e00';
const lightBackground = '#e3f2fd';
const darkBackground = '#42484d';

function FillBar(props) {
    let width = props.sum / props.amount * 100
    let color = green;
    if (props.sum < 0 && props.amount > 0) {
        width = 100;          // Bar width is meaningless if sign(sum) != sign(amount)
        color = red;          // Expected net positive, currently net negative.

    } else if (props.sum > 0 && props.amount < 0) {
        width = 0;            // Expected net negative, currently net positive.
    } else if (width > 100) { // Category is beyond budget (but has expected sign)
        width = 100;
        if (props.sum > 0) {  // Good for income
            color = green;
        } else {              // Bad for expenses
            color = red;
        }
    } else if (props.amount < 0 && Math.abs(props.sum / props.amount) > ((1 / 12) + (props.now / 100))) {
        color = yellow;       // Expense target is running beyond target
    }


    let background;
    const isDarkModeEnabled = useMediaQuery('(prefers-color-scheme: dark)');
    background = (isDarkModeEnabled ? darkBackground : lightBackground);

    // Size the bars according to the amount of money represented in them.
    let height;
    if (props.amount > 0) {
        height = "28px";
    } else {
        height = Math.abs(props.amount * (56 / 15000)) + "px";
    }

    let style = {
        'width': String(width) + "%",
        'backgroundColor': color,
        'height': height
    }

    let nowStyle = {
        'left': String(props.now) + "%",
        'width': "2px",
        'position': "relative",
        'backgroundColor': 'black',
        'top': "-" + height,
        'height': height
    }

    return (
        <div className="fillBar" style={{
            "backgroundColor": background,
            "height": height
        }}>
            <div className="filled" style={style}></div>
            <div className="nowLine" style={nowStyle}></div>
        </div>
    );
}

export default FillBar;
