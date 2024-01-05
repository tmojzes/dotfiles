return {
	{
		"nvimdev/dashboard-nvim",
		event = "VimEnter",
		opts = function(_, opts)
			local logo = [[
             ████████╗███╗   ███╗ ██████╗      ██╗███████╗███████╗███████╗   ███╗   ██╗██╗   ██╗██╗███╗   ███╗
             ╚══██╔══╝████╗ ████║██╔═══██╗     ██║╚══███╔╝██╔════╝██╔════╝   ████╗  ██║██║   ██║██║████╗ ████║
                ██║   ██╔████╔██║██║   ██║     ██║  ███╔╝ █████╗  ███████╗   ██╔██╗ ██║██║   ██║██║██╔████╔██║
                ██║   ██║╚██╔╝██║██║   ██║██   ██║ ███╔╝  ██╔══╝  ╚════██║   ██║╚██╗██║╚██╗ ██╔╝██║██║╚██╔╝██║
                ██║   ██║ ╚═╝ ██║╚██████╔╝╚█████╔╝███████╗███████╗███████║██╗██║ ╚████║ ╚████╔╝ ██║██║ ╚═╝ ██║
                ╚═╝   ╚═╝     ╚═╝ ╚═════╝  ╚════╝ ╚══════╝╚══════╝╚══════╝╚═╝╚═╝  ╚═══╝  ╚═══╝  ╚═╝╚═╝     ╚═╝
      ]]

			logo = string.rep("\n", 8) .. logo .. "\n\n"
			opts.config.header = vim.split(logo, "\n")
		end,
	},

	-- disable indentscope animation
	{
		"echasnovski/mini.indentscope",
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
}
