package views

const Styles = `
body {
    font: 16px/1.3em Georgia;
    margin: 2rem;
}

header {
    margin: 1rem 0 2rem;
}

.nowplaying {
    margin-top: 1em;
}

section {
    margin: 1rem 0 2rem;
    border-top: 1px dotted;
}

h1 {
    font-size: 1.5rem;
}

h1 a {
    color: black;
    text-decoration: none;
}

h2 {
    font-size: 1rem;
    font-variant: small-caps;
}

.tabs-choice {
    margin-bottom: .5rem;
}

.tabs-choice h3 {
    margin: 0;
    color: #666;
}

.tabs-choice h3.selected {
    text-decoration: underline;
    color: black;
}

.tabs-content {
    margin: 0;
}

ol {
    list-style: none;
    padding-left: 0;
}

.tabs-choice {
    display: flex;
}

.tabs-choice h3 {
    cursor: pointer;
    font-size: 1rem;
}

.tabs-choice h3 + h3 {
    margin-left: 1em;
}

.tabs-content {
    list-style: none;
    padding: 0;
}

.tabs-content .hide {
    display: none;
}

section {
    max-width: 40rem;
}

.more {
    display: block;
    text-align: right;
    padding-top: .5rem;
    height: .25rem;
    text-decoration: underline;
    color: black;
}

table {
    border-collapse: collapse;
    width: 100%;
    table-layout: fixed;
}

td {
    padding: 0;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
}

tr td:last-child {
    width: 180px;
}

.artist {
    font-weight: bold;
}

.artist, .album, .track {
    color: black;
    text-decoration: none;
}

.artist:hover, .album:hover, .track:hover {
    border-bottom: 1px solid #ccc;
}

.artist + .track {
    margin-left: .2rem;
}

p .artist, p .album, p .track {
    border-bottom: 1px solid #ccc;
}

a[href^='/listen'] {
    color: black;
}
time, .count {
    float: right;
    font-style: italic;
}
time.unright {
    float: inherit;
}

.hero {
    margin: 0;
}

.hero h1 {
    position: absolute;
    margin-left: 3rem;
    margin-top: 3rem;
}

.graph {
    display: flex;
    flex-wrap: nowrap;
    flex-direction: row;
    flex-flow: space-between;
    align-items: flex-end;
    list-style: none;
    height: 6rem;
    padding-left: 0;
}

.graph li {
    width: 1rem;
    display: block;
    background: #ccc;
}

@media screen and (max-width: 30em) {
    table tr td:last-child { display: none; }
}
`
