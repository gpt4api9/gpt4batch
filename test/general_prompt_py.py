import os
import fnmatch
import json


# 递归搜索目录以找到匹配的图片文件
def find_images(directory, patterns=None):
    if patterns is None:
        patterns = ['*.png', '*.jpg', '*.jpeg', '*.gif', '*.bmp', '*.webp']
    matches = []
    for root, dirnames, filenames in os.walk(directory):
        for pattern in patterns:
            for filename in fnmatch.filter(filenames, pattern):
                matches.append(os.path.join(root, filename))
    return matches


def format_image_path(image_path):
    parent_dir = os.path.basename(os.path.dirname(image_path))
    filename = os.path.basename(image_path)
    return os.path.join(parent_dir, filename)


# 读取所有图片和prompt分行写入到out.jsonl文件
def read_images_to_out_js(directory, prompt):
    images = find_images(directory)
    with open("out.jsonl", "w", encoding="utf-8") as f:
        for idx, image in enumerate(images):
            formatted_image_path = format_image_path(image)
            data = {
                "id": f"in-{idx}",
                "asks": [
                    {
                        "id": f"ask-{idx}",
                        "content": prompt,
                        "images": [formatted_image_path],
                    }
                ],
                "extra": formatted_image_path,
            }
            f.write(json.dumps(data, ensure_ascii=False) + "\n")


if __name__ == '__main__':
    read_images_to_out_js("/Users/img", "Tell me!")
