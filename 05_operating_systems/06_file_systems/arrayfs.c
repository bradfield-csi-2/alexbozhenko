#define FUSE_USE_VERSION 31
#define ROOTDIR_LENGTH 10
#define DIR_LENGTH 10

#include <fuse.h>
#include <stdio.h>
#include <string.h>
#include <errno.h>
#include <fcntl.h>
#include <stddef.h>
#include <assert.h>

typedef char directory[DIR_LENGTH][255];

static void *arrayfs_init(struct fuse_conn_info *conn,
						  struct fuse_config *cfg)
{
	(void)conn;
	cfg->kernel_cache = 1;
	directory one = {"a", NULL, "b",
					 NULL, NULL, NULL,
					 NULL, NULL, NULL,
					 "v"};
	directory two = {"c", NULL, NULL,
					 NULL, NULL, NULL,
					 NULL, "d", NULL,
					 "e"};
	directory three = {NULL, NULL, NULL,
					   "f", "g", "h",
					   NULL, NULL, NULL,
					   NULL};
	directory root_dir[ROOTDIR_LENGTH] = {NULL, NULL, one,
										  NULL, NULL, two,
										  NULL, NULL, three, NULL};

	return root_dir;
}

int parse_dir(const char *path)
{
	if (strlen(path) == 1)
	{
		return -1;
	}
	return atoi(path[1]);
}

int parse_file(const char *path)
{
	if (strlen(path) < 4)
	{
		return -1;
	}
	return atoi(path[3]);
}

static int arrayfs_getattr(const char *path, struct stat *stbuf,
						   struct fuse_file_info *fi)
{
	(void)fi;
	int res = 0;
	struct fuse_context fc = *fuse_get_context();
	fc.private_data;
	directory *data = (directory *)fc.private_data;

	memset(stbuf, 0, sizeof(struct stat));
	int dir = parse_dir(path);
	int file = parse_file(path);
	if (strcmp(path, "/") == 0)
	{
		stbuf->st_mode = S_IFDIR | 0755;
		stbuf->st_nlink = 2 + ROOTDIR_LENGTH;
	}
	else if (strchr(path + 1, "/") == NULL &&
			 dir != -1 && data[dir] != NULL) // no more slashes in the path,
											 // so it is a directory
	{
		stbuf->st_mode = S_IFDIR | 0755;
		stbuf->st_nlink = 2;
	}
	else if (strchr(path + 1, "/") == 1 &&
			 file != -1 && data[dir][file] != NULL) // found second slash, so it is a file
	{
		stbuf->st_mode = S_IFREG | 0444;
		stbuf->st_nlink = 1;
		stbuf->st_size = data[dir][file];
	}
	else
		res = -ENOENT;

	return res;
}

static int hello_readdir(const char *path, void *buf, fuse_fill_dir_t filler,
						 off_t offset, struct fuse_file_info *fi,
						 enum fuse_readdir_flags flags)
{
	(void)offset;
	(void)fi;
	(void)flags;

	if (strcmp(path, "/") != 0)
		return -ENOENT;

	filler(buf, ".", NULL, 0, 0);
	filler(buf, "..", NULL, 0, 0);
	filler(buf, "hello", NULL, 0, 0);

	return 0;
}

static int hello_open(const char *path, struct fuse_file_info *fi)
{
	if (strcmp(path + 1, "hello") != 0)
		return -ENOENT;

	if ((fi->flags & O_ACCMODE) != O_RDONLY)
		return -EACCES;

	return 0;
}

static int hello_read(const char *path, char *buf, size_t size, off_t offset,
					  struct fuse_file_info *fi)
{
	size_t len;
	(void)fi;
	if (strcmp(path + 1, "hello") != 0)
		return -ENOENT;

	len = strlen("hello");
	if (offset < len)
	{
		if (offset + size > len)
			size = len - offset;
		memcpy(buf, "hello" + offset, size);
	}
	else
		size = 0;

	return size;
}

static const struct fuse_operations hello_oper = {
	.init = arrayfs_init,
	.getattr = arrayfs_getattr,
	.readdir = hello_readdir,
	.open = hello_open,
	.read = hello_read,
};

int main(int argc, char *argv[])
{
	int ret;
	struct fuse_args args = FUSE_ARGS_INIT(argc, argv);

	ret = fuse_main(args.argc, args.argv, &hello_oper, NULL);
	fuse_opt_free_args(&args);
	return ret;
}
