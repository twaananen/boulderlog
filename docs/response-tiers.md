# Grade-relative boulder responses

Implemented behavior and editorial guide. The application embeds the approved catalogue from `services/response-messages.json`; the standalone preview is in `docs/response-preview.html`.

## Behavior

After a successful log save, choose one of six response tiers using the outcome and the distance below the user's highest recorded topped grade. Each step on the existing grade picker counts separately, as agreed:

`4, 4+, 5, 5+, 6A, 6A+, 6B, 6B+, 6C, 6C+, 7A, 7A+, 7B, 7B+, 7C, 7C+, 8A, 8A+, 8B, 8B+, 8C`

Use this explicit order, never string sorting or whole-number subtraction. Let `gap = index(previousBestTop) - index(loggedGrade)`.

| Grade relative to previous best top | Topped: difficulty 1–4 | Not topped: difficulty 5–8 |
| --- | --- | --- |
| Above best, at best, or 1–2 steps below (`gap <= 2`) | `victory`: enthusiastic celebration | `pep_talk`: encouragement for a real challenge |
| 3–5 steps below (`3 <= gap <= 5`) | `mildly_impressed`: a modest nod of approval | `light_roast`: sharp, playful teasing |
| 6 or more steps below (`gap >= 6`) | `unimpressed`: dry acknowledgement of an easy-for-you send | `full_roast`: a much harsher roast of the attempt |

For a previous best of **7A+**: 6C+ and higher use victory / pep talk; 6B through 6C use mildly impressed / light roast; **6A+ and lower use unimpressed / full roast**. In particular, failing 6A+ is six picker steps below 7A+ and failing 6A is seven: both get a full roast. A failed attempt above the best is always a pep talk.

## Baseline and edge cases

- **Best means all-time best among the authenticated user's existing, non-deleted, valid topped logs.** A top requires difficulty 1–4. Failed attempts, flash status and the `NewRoute` flag do not change the baseline. Imported logs count if their grades match the picker and their difficulty is valid.
- **Read the baseline before saving the new log.** A first top or new personal best gets a victory line; subsequent logs use the updated history. A baseline read error must not prevent saving the climb: use the safe fallback below and log the lookup error.
- **No prior valid top:** topped climbs get `victory`; unsuccessful attempts get `pep_talk`. Do not infer ability from failed attempts.
- **Unrecognized historical grades or invalid historical difficulty:** skip them when finding the best. Match trimmed, uppercased grade strings to the explicit picker list; do not silently map V grades or other systems onto it. All-invalid history is equivalent to no baseline.
- **New input:** accept only a supported picker grade and difficulty 1–8. Invalid values are validation errors before a save, not a way to reach a roast tier. Use the same grade validation on edits so new history stays comparable.
- **If response context is unavailable:** use victory for a valid top and pep talk for a valid unsuccessful attempt. A missing or failed baseline lookup must never produce a roast.
- **Edits and deletes:** recompute the best from current history on the next new log. Do not send a new celebration or roast merely for editing/deleting an entry. No cached best-grade field is needed.
- **Flash / repeats:** these do not override the selected tier in this version. All topped difficulties use the same grade-relative bands; all not-topped difficulties do likewise.
- **Ordinary home navigation:** keep the current behavior of showing no post-log response.

Known limit: an all-time best can overestimate current ability after a long break or during a recovery period. This proposal follows the requested all-time reference; a recent-best window or tone preference is deferred rather than added implicitly. Grades and styles vary between gyms, so the roast is a joke about the recorded comparison, not a reliable assessment of effort.

## Content and tone

`services/response-messages.json` contains **170 messages per tier, 1,020 total**. Every theme has at least twice its original five lines in every tier. Five especially rich themes have 20 each, giving them more room and naturally more chances to appear.

| Theme key | Lines per tier | Coverage |
| --- | --- | --- |
| `general` | 10 | Direct celebration, acknowledgement or teasing |
| `climbing_humour` | 10 | Chalk, gravity, holds, beta, sending and climbing wordplay |
| `learning_persistence` | 10 | Practice, experiments, experience and progress |
| `mythology` | 20 | Greek, Norse, Egyptian, Chinese, Japanese and West African figures |
| `middle_earth` | 20 | Tolkien characters, places and journeys |
| `fantasy_magic` | 10 | Conan, Harry Potter, Merlin, Discworld and Howl's Moving Castle |
| `fantasy_worlds` | 20 | Witcher, Game of Thrones, Elder Scrolls, Zelda and God of War |
| `superheroes` | 10 | Marvel, DC and Masters of the Universe |
| `science_fiction` | 10 | Star Wars, Matrix, Mass Effect, Halo, Alien and Star Trek |
| `gaming` | 20 | Assassin's Creed, Minecraft, Mario, Sonic, Souls / Elden Ring, Hollow Knight, Doom, Red Dead Redemption, Baldur's Gate 3, Disco Elysium and Hades |
| `anime` | 20 | Attack on Titan, Dragon Ball, Naruto, One Piece, One-Punch Man, JoJo, Death Note, My Hero Academia and anime tropes |
| `portal_science` | 10 | Portal, GLaDOS, Wheatley, Companion Cube, Cave Johnson, Chell and puzzle science |

Every line was rewritten again in November 2026 to be funny first: invented in-character quotes, riffs on famous lines, short punchlines and absurd escalation. Each of the 1,020 lines was graded individually by a separate reviewer model on funniness and theme fit and rewritten until it scored at least 8/10 on both. The full-roast pool is the harshest, with the most profanity; light roasts allow a little. Humour should have a specific image or situation rather than just replacing one name in a stock sentence. The bank is grouped by theme for editorial review; theme keys are not an extra user-facing choice.

### Recurring characters

The preview has a character comparison view. Its roster contains **38 named characters or figures**, each with at least one line in all six tiers within their theme. Many other characters appear in fewer tiers; the comparison view shows "—" for any tier without a line. This is editorial coverage, not a new application preference or a requirement to repeat the same character to a user.

Examples of complete coverage:

- **Mythology:** Thor, Loki, Zeus, Sisyphus, Odin and Hercules.
- **Middle-earth:** Gandalf, Gollum, Gimli, Aragorn and Bilbo.
- **Fantasy magic:** Conan, Hermione, Ron and Granny Weatherwax.
- **Fantasy worlds:** Geralt, Kratos, Dandelion, Tormund, Atreus and Mimir.
- **Superheroes:** Deadpool, Captain America, Tony, He-Man and Skeletor.
- **Science fiction:** Han.
- **Gaming:** Mario and Kim Kitsuragi.
- **Anime:** Levi, Hange, Goku, Vegeta, Kakashi and Gai.
- **Portal:** GLaDOS, Wheatley and Cave Johnson.
- The roster is embedded in the standalone preview and checked against the actual bank.

### Tone

- Victory should feel earned, without claiming a personal best, flash, session streak or particular move that the selector has not verified.
- Pep talks acknowledge the challenge without assuming the user nearly topped. They work for all four unsuccessful difficulty ratings.
- Mild approval is noticeably warmer than indifference. Unimpressed lines acknowledge completion without treating it as failure or mocking the decision to warm up.
- Light roast is sharp, playful teasing. Full roast should be much harsher: exaggerated humiliation, ruined reputations, absurd comparisons and occasional profanity. Both target this attempt, the claimed climbing credentials and the gap to the recorded best, not appearance, identity or personal worth.
- Stay in the reference's world and in character. Gandalf talks like Gandalf (riddles, wizard timing, "you shall not pass"), GLaDOS gives passive-aggressive test results, Kim Kitsuragi writes one word in his notebook. Invented quotes are encouraged: put new in-voice lines in characters' mouths, riff on their famous lines and land an unexpected punchline. Long verbatim quotes are not. Keep the lore right; a joke that contradicts the source reads as a mistake.
- Never put stock app praise in a character's mouth: no "X approves!", "X would be proud!", "What a send!" or "nice boulder". A line needs a specific image or twist from that world mapped onto the climb, not a name swapped into a generic sentence.
- Keep lines short, readable and enjoyable without knowing the reference well. Deadpan beats shouting; use few exclamation marks. No advice to ignore pain, skip rest or take unsafe risks.

## Selection and application changes

1. Validate the new log, obtain the user's previous best using the existing `Database.GetBoulderLogs(username)` history method, and save the log. Handle a history lookup failure with the neutral fallback; handle a save failure with the existing error path and no response.
2. Compute the tier once in the service layer. For a known baseline use the table above; otherwise use the fallback tiers.
3. Flatten the selected tier's 12 theme arrays and choose uniformly among its 170 messages. There is no second theme draw or separate weighting configuration. Each 20-line theme has twice the chance of each 10-line theme; the five rich themes collectively supply 100 of 170 options. Repeats are allowed, matching the current random selection behavior.
4. Pass the selected message into `LogSummary`; the template renders escaped text and retains the existing topped / not-topped visual distinction. Keep normal home renders free of a response.

Implementation:

- `services/log_service.go`: user-scoped baseline lookup and tier decision, reusing the existing database interface and topped definition.
- `services/responses.go` embeds `services/response-messages.json` with `go:embed` and loads it once. The catalogue is part of the binary; no runtime file path or additional dependency is needed.
- `components/boulder_grade_selection.templ`: use the same ordered grade list as response comparison so the UI and ranking cannot drift. Validation and comparison reuse that list; chart sorting is outside this feature.
- `handlers/log_handler.go`: pass response context through the save path and validate supported grades / difficulty on create and edit.
- `components/log_summary.templ`, `components/home.templ`, `handlers/home_handler.go`: update callers for the selected-message input, preserving the existing no-response home behavior. Remove the two old message functions once the new bank is wired in. Regenerate affected tracked `*_templ.go` files.
- Service / handler / component tests and existing database mocks only where their calls or signatures change. No database schema, migration, settings screen, external API or theme preference is required.

History scans are the smallest use of the existing API. If per-user history size makes that lookup materially slow, replace it with a database-side maximum over the same explicit rank order. Do not fetch other users' logs or introduce a stored personal-best field for this version.

## Acceptance checks for implementation

- At a 7A+ best, test gaps `-1, 0, 1, 2, 3, 4, 5, 6` against both outcomes: exactly the tiers in the table. Explicitly verify failed 6A+ and 6A both produce full roast, failed 6B and 6C produce light roast, and failed 6C+ produces pep talk.
- Test grade boundaries `5+ → 6A` and `6C+ → 7A`, and the low/high ends `4` and `8C`. A plus or letter change is one picker step.
- Test all difficulty values 1–8: 1–4 are tops, 5–8 are unsuccessful. Reject 0, 9, malformed or unsupported new grades before saving.
- A user's highest failed grade must not become their baseline; another user's top must not influence it. Ignore deleted and invalid historical records.
- Test no history, only failures, all-invalid history, first top, a new high, baseline lookup failure and save failure. Missing baseline yields no roast; a failed save yields no success response.
- After editing/deleting the previous highest top, the next log uses the remaining best. Editing itself produces no response.
- Assert the bank has exactly the six expected tiers, 170 lines per tier, all 12 themes, 20 lines per rich theme and 10 per other theme, no empty lines or duplicate messages, and no message over 180 characters. Check all 38 comparison characters appear in all six tiers in their own theme. Check actual message selection returns a member of its chosen tier, without flaky distribution assertions.
- Render one message from each tier and verify text escaping and the normal home state. Exercise the HTMX post-log flow and direct-page rendering.

## Validation

The feature is integrated into the log submission flow. Automated service, handler and rendering tests cover the tier boundaries, baseline filtering and failure fallback, validation, save ordering, both HTMX and full-page responses, and text escaping. The catalogue has count, uniqueness and length checks. The preview also verifies all 38 roster characters across six tiers.

Run the Go suite with `go test -race . ./components ./db ./handlers ./middleware ./models ./services ./utils`, then run `node --test tests/charts.test.cjs`. Use explicit package paths because local development database files under `data/postgres` may prevent recursive `./...` package discovery.
