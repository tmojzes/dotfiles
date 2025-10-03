return {
	{
		"folke/snacks.nvim",
		opts = {
			dashboard = {
				preset = {
					header = [[
             ████████╗███╗   ███╗ ██████╗      ██╗███████╗███████╗███████╗   ███╗   ██╗██╗   ██╗██╗███╗   ███╗
             ╚══██╔══╝████╗ ████║██╔═══██╗     ██║╚══███╔╝██╔════╝██╔════╝   ████╗  ██║██║   ██║██║████╗ ████║
                ██║   ██╔████╔██║██║   ██║     ██║  ███╔╝ █████╗  ███████╗   ██╔██╗ ██║██║   ██║██║██╔████╔██║
                ██║   ██║╚██╔╝██║██║   ██║██   ██║ ███╔╝  ██╔══╝  ╚════██║   ██║╚██╗██║╚██╗ ██╔╝██║██║╚██╔╝██║
                ██║   ██║ ╚═╝ ██║╚██████╔╝╚█████╔╝███████╗███████╗███████║██╗██║ ╚████║ ╚████╔╝ ██║██║ ╚═╝ ██║
                ╚═╝   ╚═╝     ╚═╝ ╚═════╝  ╚════╝ ╚══════╝╚══════╝╚══════╝╚═╝╚═╝  ╚═══╝  ╚═══╝  ╚═╝╚═╝     ╚═╝
]],
				},
			},
		},
	},

	-- disable indentscope animation
	{
		"nvim-mini/mini.indentscope",
		opts = {
			draw = { animation = require("mini.indentscope").gen_animation.none() },
		},
	},

	-- show hidden files in neotree
	{
		"nvim-neo-tree/neo-tree.nvim",
		opts = {
			filesystem = {
				filtered_items = {
					visible = false,
					show_hidden_count = true,
					hide_dotfiles = true,
					hide_gitignored = true,
					hide_by_name = {
						".git",
						".DS_Store",
					},
					never_show = {},
				},
			},
		},
	},
	{
		"Bekaboo/dropbar.nvim",
		dependencies = {
			-- optional, but required for fuzzy finder support
			"nvim-telescope/telescope-fzf-native.nvim",
			-- optional, but required if you want to see icons for different filetypes
			"nvim-tree/nvim-web-devicons",
		},
	},
}
