package main

import (
	"fmt"
	"io/ioutil"
)

const FLODER_DEPTH = 2 // traverse depth

var floders = map[int]map[string][]string{} // depth[floder]floders
var files = map[string][]string{}           // floder[files]
var filter_folders = map[string]bool{       // filter folders in map
	"img":  true,
	".git": true,
}
var filter_files = map[string]bool{
	".DS_Store":     true,
	"generate.go":   true,
	"config.json":   true,
	"index.html":    true,
	"index.md":      true,
	"navigation.md": true,
	"wiki_index.md": true,
	// "*.sql": true,
	// "*.pdf": true,
} // filter files in map

var idx_md = "wiki_index.md"        // index md file name
var navigation_md = "navigation.md" // nav md file name

var need_chose_theme = false // is need to chose a wiki theme switch in nav

func main() {
	//traverse all file

	floder_depth := map[string]int{}

	floders, files, _ = get_all_files(".", floders, files, 0, floder_depth, filter_folders, filter_files)

	fmt.Printf("floders: %v\n", floders)
	fmt.Println("==============================")
	fmt.Printf("files: %v", files)

	generate_nav(navigation_md, idx_md, floders, files)
	generate_index(idx_md, floders, files)
}

func generate_nav(navigation_md, index_md string, floders map[int]map[string][]string, files map[string][]string) {
	text := "[gimmick:theme](cerulean)\n"
	text += "# Gavin Wiki" // wiki name
	for depth, child_floders := range floders {
		if depth == FLODER_DEPTH {

			for f, ch_fs := range child_floders {
				var f_name string
				if len(f) >= 2 {
					f_name = f[2:]
				} else {
					f_name = f
				}

				text += "\n"

				if len(ch_fs) == 0 {
					text += "[" + f_name + "](" + f + "/" + index_md + ")\n"
				} else {
					text += "[" + f_name + "]()\n"
					text += "\n"
					for _, fs := range ch_fs {
						var fs_name string
						if len(fs) >= 2 {
							fs_name = fs[2:]
						} else {
							fs_name = fs
						}
						text += "  * [" + fs_name + "](" + fs + "/" + index_md + ")\n"
					}
				}
				// file idx
				if t_fs, ok := files[f]; ok {
					for _, t_f := range t_fs {
						text += "  * [" + t_f + "](" + f + "/" + t_f + ")\n"
					}
				}
			}

		}

	}
	if need_chose_theme {
		text += "\n[gimmick:themechooser](Choose theme)"
	}

	if err := write(navigation_md, text); err != nil {
		fmt.Printf("generate nav err, err=%v\n", err)
	}
}

func generate_index(index_md string, floders map[int]map[string][]string, files map[string][]string) {

	idx_text := "Welcome\n=======\nHave fun\n----------\n这里是 Gavin 的 wiki 空间\n"

	for depth, child_floders := range floders {
		if depth == FLODER_DEPTH {

			for f, ch_fs := range child_floders {
				var f_name string
				if len(f) >= 2 {
					f_name = f[2:]
				} else {
					f_name = f
				}

				idx_text += "\n"

				idx_text += "[" + f_name + "](" + f + "/" + index_md + ")\n\n"
				if len(ch_fs) > 0 {
					if t_fs, ok := files[f]; ok {
						for _, t_f := range t_fs {
							idx_text += "  * [" + t_f + "](" + f + "/" + t_f + ")\n"
						}
					}

					for _, fs := range ch_fs {
						var fs_name string
						if len(fs) >= 2 {
							fs_name = fs[2:]
						} else {
							fs_name = fs
						}
						idx_text += "  * [" + fs_name + "](" + fs + "/" + index_md + ")\n"
					}
				}
			}

		}

	}

	if err := write("index.md", idx_text); err != nil {
		fmt.Printf("generate idx err, err=%v\n", err)
	}

	for depth, child_floders := range floders {
		if depth == FLODER_DEPTH {
			for f, ch_fs := range child_floders {

				var f_name string
				if len(f) >= 2 {
					f_name = f[2:]
				} else {
					f_name = f
				}

				text := "### " + f_name + "\n"
				text += "----------\n"

				f_path := f + "/" + index_md

				if t_fs, ok := files[f]; ok {
					for _, t_f := range t_fs {
						text += "[" + t_f + "](" + t_f + ")\n"
					}
				}

				if err := write(f_path, text); err != nil {
					fmt.Printf("generate index err, path=%v, err=%v\n", f_path, err)
				}

				if len(ch_fs) >= 0 {

					for _, fs := range ch_fs {
						var fs_name string
						if len(fs) >= 2 {
							fs_name = fs[2:]
						} else {
							fs_name = fs
						}

						text := "### " + fs_name + "\n"
						text += "----------\n"

						f_path := fs + "/" + index_md

						if t_fs, ok := files[fs]; ok {
							for _, t_f := range t_fs {
								text += "[" + t_f + "](" + t_f + ")\n"
							}
						}

						if err := write(f_path, text); err != nil {
							fmt.Printf("generate index err, path=%v, err=%v\n", f_path, err)
						}

					}
				}
			}

		}

	}
}

func get_all_files(pathname string, floders map[int]map[string][]string, files map[string][]string, parent_depth int, floder_depth map[string]int, filter_folders map[string]bool, filter_files map[string]bool) (
	map[int]map[string][]string, map[string][]string, error) {

	if d, ok := floder_depth[pathname]; ok && d >= FLODER_DEPTH {
		return floders, files, nil
	}

	rd, err := ioutil.ReadDir(pathname)
	if err != nil {
		fmt.Println("read dir fail:", err)
		return floders, files, err
	}

	if _, ok := floder_depth[pathname]; !ok {
		floder_depth[pathname] = parent_depth + 1
	} else {
		floder_depth[pathname]++
	}

	depth := floder_depth[pathname]
	if _, ok := floders[depth]; !ok {
		floders[depth] = make(map[string][]string, 0)
	}
	if _, ok := floders[depth][pathname]; !ok {
		floders[depth][pathname] = make([]string, 0)
	}
	if _, ok := files[pathname]; !ok {
		files[pathname] = make([]string, 0)
	}

	for _, fi := range rd {
		if fi.IsDir() {

			if _, ok := filter_folders[fi.Name()]; ok {
				continue
			}

			fullDir := pathname + "/" + fi.Name()
			floders[depth][pathname] = append(floders[depth][pathname], fullDir)

			floders, files, err = get_all_files(fullDir, floders, files, depth, floder_depth, filter_folders, filter_files)
			if err != nil {
				fmt.Println("read dir fail:", err)
				return floders, files, err
			}
		} else {

			if _, ok := filter_files[fi.Name()]; ok {
				continue
			}

			files[pathname] = append(files[pathname], fi.Name())
		}
	}
	return floders, files, nil
}

func write(path string, text string) (err error) {
	var d1 = []byte(text)
	err = ioutil.WriteFile(path, d1, 0666)
	return err
}
