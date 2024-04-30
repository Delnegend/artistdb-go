# artistdb-go

## Developement
- Install [Go](https://go.dev/dl/), [air](https://github.com/cosmtrek/air#installation), [nodejs](https://nodejs.org/en) and [pnpm](https://pnpm.io/)
- `pnpm i && go mod tidy`
- `pnpm dev` to start the server, listening on port `8080`
- `pnpm format` to format the `html`, `css` and `go` files

## Environment variables

| Name | Description | Default |
| --- | --- | --- |
| `PORT` | Port to listen on | `8080` |
| `IN_FILE` | Path to the input file | `artists.txt` |
| `OUT_DIR` | Path to the output directory | `artists` |
| `AVATAR_DIR` | Path to the avatar directory | `avatar` |
| `FORMAT_AND_EXIT` | Sort the artists alphabetically by username, then overwrite the input file with the sorted data | `false` |
| `FALLBACK_AVATAR` | Path to the fallback avatar for unavatar, must be accessible from the web | `https://via.placeholder.com/150` |

## artists.txt file structure
```
username[,displayName,avatar,...alias]
[*,]social[,description]
...
```

- All username and alias must be unique
- Avatar has 2 format: `username@social` or `/path/to/image.format`
    > Create a directory at the root of the project, place the image in it and use "/<filename>.<fileformat>", "./avatar" will be automatically added to the front of the path.
- `*` socials will have a more highlighted format on the frontend.
- `social` has 2 format
    - `username@social`: description is optional if provided, the description will have `<social_name> |` prepended to it
        > For example, `foo@instagram,Personal` -> description = `Instagram | Personal`
    - `//example.com/username`: description is required

### Example
```
paul,Paul Something,paul@twitter,paulsomething
*//example.com/paul,Paul's website
paul@twitter,Life
paulart@twitter,Art account

john,John Doe,john@twitter,johndoe
...
```

## Output file structure
### Main file
```
displayName,avatar
socialLink,description
```
- Has username as the filename
- `avatar` has 2 format:
    - `social/username`
    - `/username.format`

### Alias file

```
@username
```

- Points to the main file