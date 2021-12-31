import * as React from 'react';

export default function AccountIcon(props) {
    /* This is just a table of my personal asset account IDs. If you're also
    using lychnos, you'll have to modify this to match your needs  :)
    TODO(davidschlachter): generate this list dynamically based on the contents
    of public/account_icons
    */
    let icons = {
        "0": "0.png", // default icon
        "1": "1.png",
        "3": "3.png",
        "5": "5.png",
        "11": "11.png",
        "28": "28.png",
        "53": "53.png",
        "71": "71.png",
        "161": "161.png",
        "469": "469.png",
        "470": "470.png",
        "483": "483.png",
    };
    let icon = icons["0"];
    if (props.account_id in icons) {
        icon = icons[props.account_id];
    }
    let iconPath = "/app/account_icons/" + icon;
    return (
        <img style={{ "height": "32px" }} src={iconPath} alt="account icon" />
    );
}